package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"route256/cart/pkg/logger"
	"route256/loms/internal/handler"
	"route256/loms/internal/infra/config"
	"route256/loms/internal/infra/grpc/interceptor"
	"route256/loms/internal/infra/http/middleware"
	"route256/loms/internal/infra/repository/postgres"
	"route256/loms/internal/service"
	"route256/loms/pkg/api/orders/v1"
	"route256/loms/pkg/api/stocks/v1"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq" // Import postgres driver
	"github.com/pressly/goose/v3"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// App создает компоненты для сервиса loms
type App struct {
	Config       *config.Config
	grpcServer   *grpc.Server
	grpcGWServer *http.Server
}

// NewApp конструктор главного приложения.
func NewApp(configPath string) (*App, error) {
	c, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("config.LoadConfig: %w", err)
	}

	app := &App{Config: c}
	app.grpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.Logging,
			interceptor.Validate,
		),
	)

	reflection.Register(app.grpcServer)

	postgresMasterDSN := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", app.Config.MasterDB.User, app.Config.MasterDB.Password,
		app.Config.MasterDB.Host, app.Config.MasterDB.Port, app.Config.MasterDB.DBName)

	postgresReplicaDSN := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", app.Config.ReplicaDB.User, app.Config.ReplicaDB.Password,
		app.Config.ReplicaDB.Host, app.Config.ReplicaDB.Port, app.Config.ReplicaDB.DBName)

	err = applyMigrations(postgresMasterDSN, app.Config.Server.MigrationsPath)
	if err != nil {
		logger.Warning(fmt.Sprintf("Skip migrations applying with error: %s", err.Error()))
	}

	poolManager, err := postgres.NewRRPoolManager(context.TODO(), postgresMasterDSN, []string{postgresReplicaDSN})
	if err != nil {
		return nil, fmt.Errorf("postgres.NewRRPoolManager: %w", err)
	}

	txManager := postgres.NewPgTxManager(poolManager)
	repositoryfactory := postgres.NewRepositoryFactory(poolManager)

	stockService := service.NewStockService(repositoryfactory, txManager)
	orderService := service.NewOrderService(stockService, repositoryfactory, txManager)

	stocksHandler := handler.NewStockServerGRPC(stockService)
	ordersHandler := handler.NewOrderServerGRPC(orderService)

	stocks.RegisterStockServiceV1Server(app.grpcServer, stocksHandler)
	orders.RegisterOrderServiceV1Server(app.grpcServer, ordersHandler)

	err = LoadStocks(stockService, time.Duration(app.Config.Server.LoadStocksDataTimeout)*time.Second)
	if err != nil {
		logger.Error(fmt.Sprintf("LoadStocks: %s", err))
	}

	return app, nil
}

func applyMigrations(dsn, migrationsFolder string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(db, migrationsFolder); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

// ListenAndServeGRPC запускает gRPC-сервер приложения.
func (a *App) ListenAndServeGRPC() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", a.Config.Server.GRPCPort))
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}

	logger.Info(fmt.Sprintf("Loms service listening gRPC at port %s", a.Config.Server.GRPCPort))

	return a.grpcServer.Serve(listener)
}

// ListenAndServeGRPCGateway запускает gRPC-gateway для gRPC-сервера приложений.
func (a *App) ListenAndServeGRPCGateway() error {
	conn, err := grpc.NewClient(
		fmt.Sprintf(":%s", a.Config.Server.GRPCPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("grpc.NewClient: %w", err)
	}

	gwMux := runtime.NewServeMux()
	ctx := context.Background()

	err = orders.RegisterOrderServiceV1Handler(ctx, gwMux, conn)
	if err != nil {
		return fmt.Errorf("orders.RegisterOrderServiceHandler: %w", err)
	}

	err = stocks.RegisterStockServiceV1Handler(ctx, gwMux, conn)
	if err != nil {
		return fmt.Errorf("stocks.RegisterStockServiceHandler: %w", err)
	}

	handler := middleware.CORSAllPass(gwMux)

	a.grpcGWServer = &http.Server{
		Addr:              fmt.Sprintf(":%s", a.Config.Server.HTTPPort),
		Handler:           handler,
		ReadHeaderTimeout: time.Second * time.Duration(a.Config.Server.GRPCGateWay.ReadHeaderTimeout),
		WriteTimeout:      time.Second * time.Duration(a.Config.Server.GRPCGateWay.WriteTimeout),
		IdleTimeout:       time.Second * time.Duration(a.Config.Server.GRPCGateWay.IdleTimeout),
	}

	logger.Info(fmt.Sprintf("Loms service listening gRPC-Gateway (REST) at port %s", a.Config.Server.HTTPPort))

	return a.grpcGWServer.ListenAndServe()
}

// Shutdown gracefully останавливает приложение.
func (a *App) Shutdown(ctx context.Context) error {
	errGroup := new(errgroup.Group)
	errGroup.Go(func() error {
		return a.grpcGWServer.Shutdown(ctx)
	})

	errGroup.Go(func() error {
		successGraceful := make(chan struct{})
		go func() {
			a.grpcServer.GracefulStop()
			close(successGraceful)
		}()

		select {
		case <-ctx.Done():
			a.grpcServer.Stop()
			return ctx.Err()
		case <-successGraceful:
			return nil
		}
	})

	return errGroup.Wait()
}
