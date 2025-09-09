package app

import (
	"fmt"
	"net"
	"route256/cart/pkg/logger"
	"route256/loms/internal/handler"
	"route256/loms/internal/infra/config"
	"route256/loms/internal/infra/repository"
	"route256/loms/internal/service"
	"route256/loms/pkg/api/orders/v1"
	"route256/loms/pkg/api/stocks/v1"

	"google.golang.org/grpc"
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
	app.grpcServer = grpc.NewServer()

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

	return app, nil
}

// ListenAndServe запускает gRPC-сервер приложения.
func (a *App) ListenAndServe() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", a.config.Server.GRPCPort))
	if err != nil {
		return nil
	}

	logger.Info(fmt.Sprintf("Loms service listening gRPC at port %s", a.config.Server.GRPCPort))

	return a.grpcServer.Serve(listener)
}
