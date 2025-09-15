package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"route256/cart/pkg/logger"
	"route256/loms/internal/handler"
	"route256/loms/internal/infra/config"
	"route256/loms/internal/infra/grpc/interceptor"
	"route256/loms/internal/infra/http/middleware"
	"route256/loms/internal/infra/repository"
	"route256/loms/internal/service"
	"route256/loms/pkg/api/orders/v1"
	"route256/loms/pkg/api/stocks/v1"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
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

	idGen := repository.NewIDGeneratorSync()
	stockRepository := repository.NewInMemoryStockRepository(10)
	orderRepository := repository.NewInMemoryOrderRepository(idGen, 10)

	stockService := service.NewStockService(stockRepository)
	orderService := service.NewOrderService(orderRepository, stockService)

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
