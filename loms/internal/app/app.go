package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof" // nolint:gosec // profiling enabled for local debugging
	"time"

	"route256/cart/pkg/logger"
	"route256/cart/pkg/myerrgroup"
	postgrespkg "route256/cart/pkg/postgres"
	"route256/cart/pkg/tracer"
	"route256/loms/internal/handler"
	"route256/loms/internal/infra/config"
	"route256/loms/internal/infra/grpc/interceptor"
	"route256/loms/internal/infra/http/middleware"
	"route256/loms/internal/infra/repository/postgres"
	"route256/loms/internal/service"
	"route256/loms/pkg/api/orders/v1"
	"route256/loms/pkg/api/stocks/v1"
	interceptorpkg "route256/loms/pkg/grpc/interceptor"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq" // Import postgres driver
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// App создает компоненты для сервиса loms
type App struct {
	Config        *config.Config
	grpcServer    *grpc.Server
	grpcGWServer  *http.Server
	tracerManager *tracer.Manager
}

// NewApp конструктор главного приложения.
func NewApp(configPath string) (*App, error) {
	c, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("config.LoadConfig: %w", err)
	}

	app := &App{Config: c}
	app.tracerManager, err = tracer.NewTracerManager(
		context.Background(),
		fmt.Sprintf("http://%s:%s", app.Config.Jaeger.Host, app.Config.Jaeger.Port),
		app.Config.Server.Tracing.ServiceName,
		app.Config.Server.Tracing.Environment,
	)
	if err != nil {
		return nil, fmt.Errorf("tracer.NewTracerManager: %w", err)
	}

	app.grpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptorpkg.NewTracing(app.tracerManager).Do,
			interceptor.Logging,
			interceptor.Metrics,
			interceptor.Validate,
		),
	)

	reflection.Register(app.grpcServer)

	postgresMasterDSN := dsnBuilder(app.Config.MasterDB.User, app.Config.MasterDB.Password,
		app.Config.MasterDB.Host, app.Config.MasterDB.Port, app.Config.MasterDB.DBName)

	postgresReplicaDSN := dsnBuilder(app.Config.ReplicaDB.User, app.Config.ReplicaDB.Password,
		app.Config.ReplicaDB.Host, app.Config.ReplicaDB.Port, app.Config.ReplicaDB.DBName)

	ctx := context.Background()
	masterPool, err := newPool(ctx, postgresMasterDSN)
	if err != nil {
		return nil, fmt.Errorf("newPool: %w", err)
	}

	replicaPools := []*pgxpool.Pool{}
	for _, dsn := range []string{postgresReplicaDSN} {
		pool, errPool := newPool(ctx, dsn)
		if errPool != nil {
			return nil, fmt.Errorf("newPool: %w", errPool)
		}
		replicaPools = append(replicaPools, pool)
	}

	poolManager, err := postgres.NewRRPoolManager(masterPool, replicaPools)
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

	return app, nil
}

func dsnBuilder(user, password, host string, port int64, dbname string) string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", user, password, host, port, dbname)
}

func newPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.ParseConfig (dsn=%s): %w", dsn, err)
	}

	config.ConnConfig.Tracer = postgrespkg.NewMetricsTracer()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.NewWithConfig (dsn=%s): %w", dsn, err)
	}

	return pool, nil
}

// ListenAndServeGRPC запускает gRPC-сервер приложения.
func (a *App) ListenAndServeGRPC() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", a.Config.Server.GRPCPort))
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}

	logger.Infow(fmt.Sprintf("Loms service listening gRPC at port %s", a.Config.Server.GRPCPort))

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

	mux := http.NewServeMux()
	mux.Handle("/", gwMux)
	mux.HandleFunc("/debug/pprof/", http.DefaultServeMux.ServeHTTP)
	mux.Handle("/metrics", promhttp.Handler())

	handler := middleware.CORSAllPass(mux)

	a.grpcGWServer = &http.Server{
		Addr:              fmt.Sprintf(":%s", a.Config.Server.HTTPPort),
		Handler:           handler,
		ReadHeaderTimeout: time.Second * time.Duration(a.Config.Server.GRPCGateWay.ReadHeaderTimeout),
		WriteTimeout:      time.Second * time.Duration(a.Config.Server.GRPCGateWay.WriteTimeout),
		IdleTimeout:       time.Second * time.Duration(a.Config.Server.GRPCGateWay.IdleTimeout),
	}

	logger.Infow(fmt.Sprintf("Loms service listening gRPC-Gateway (REST) at port %s", a.Config.Server.HTTPPort))

	return a.grpcGWServer.ListenAndServe()
}

// Shutdown gracefully останавливает приложение.
func (a *App) Shutdown(ctx context.Context) error {
	errGroup, ctx := myerrgroup.WithContext(ctx)
	errGroup.Go(func() error {
		return a.tracerManager.Stop(ctx)
	})

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
