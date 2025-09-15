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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	config     *config.Config
	grpcServer *grpc.Server
}

// NewApp конструктор главного приложения.
func NewApp(configPath string) (*App, error) {
	c, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("config.LoadConfig: %w", err)
	}

	app := &App{config: c}
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

	stocks.RegisterStockServiceServer(app.grpcServer, stocksHandler)
	orders.RegisterOrderServiceServer(app.grpcServer, ordersHandler)

	err = LoadStocks(stockService, time.Duration(app.config.Server.LoadStocksDataTimeout)*time.Second)
	if err != nil {
		logger.Error(fmt.Sprintf("LoadStocks: %s", err))
	}

	return app, nil
}

// ListenAndServeGRPC запускает gRPC-сервер приложения.
func (a *App) ListenAndServeGRPC() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", a.config.Server.GRPCPort))
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}

	logger.Info(fmt.Sprintf("Loms service listening gRPC at port %s", a.config.Server.GRPCPort))

	return a.grpcServer.Serve(listener)
}

// ListenAndServeGRPCGateway запускает gRPC-gateway для gRPC-сервера приложений.
func (a *App) ListenAndServeGRPCGateway() error {
	conn, err := grpc.NewClient(
		fmt.Sprintf(":%s", a.config.Server.GRPCPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("grpc.NewClient: %w", err)
	}

	gwMux := runtime.NewServeMux()
	ctx := context.Background()

	err = orders.RegisterOrderServiceHandler(ctx, gwMux, conn)
	if err != nil {
		return fmt.Errorf("orders.RegisterOrderServiceHandler: %w", err)
	}

	err = stocks.RegisterStockServiceHandler(ctx, gwMux, conn)
	if err != nil {
		return fmt.Errorf("stocks.RegisterStockServiceHandler: %w", err)
	}

	handler := middleware.CORSAllPass(gwMux)

	gwServer := &http.Server{
		Addr:              fmt.Sprintf(":%s", a.config.Server.HTTPPort),
		Handler:           handler,
		ReadHeaderTimeout: time.Second * time.Duration(a.config.Server.GRPCGateWay.ReadHeaderTimeout),
		WriteTimeout:      time.Second * time.Duration(a.config.Server.GRPCGateWay.WriteTimeout),
		IdleTimeout:       time.Second * time.Duration(a.config.Server.GRPCGateWay.IdleTimeout),
	}

	logger.Info(fmt.Sprintf("Loms service listening gRPC-Gateway (REST) at port %s", a.config.Server.HTTPPort))

	return gwServer.ListenAndServe()
}
