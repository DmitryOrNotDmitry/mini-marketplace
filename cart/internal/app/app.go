package app

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"route256/cart/internal/handler"
	"route256/cart/internal/infra/config"
	"route256/cart/internal/infra/http/middleware"
	"route256/cart/internal/infra/http/roundtripper"
	"route256/cart/internal/infra/repository"
	"route256/cart/internal/service"
	"route256/cart/pkg/logger"
	"route256/loms/pkg/api/orders/v1"
	"route256/loms/pkg/api/stocks/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	config *config.Config
	server http.Server
}

// NewApp конструктор главного приложения.
func NewApp(configPath string) (*App, error) {
	c, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("config.LoadConfig: %w", err)
	}

	app := &App{config: c}
	app.server.Handler, err = app.bootstrapHandlers()
	if err != nil {
		return nil, fmt.Errorf("app.bootstrapHandlers: %w", err)
	}

	return app, nil
}

// ListenAndServe запускает HTTP-сервер приложения.
func (app *App) ListenAndServe() error {
	address := fmt.Sprintf("%s:%s", app.config.Server.Host, app.config.Server.Port)

	l, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}

	logger.Info(fmt.Sprintf("Cart service listening at http://%s", address))

	return app.server.Serve(l)
}

func (app *App) bootstrapHandlers() (http.Handler, error) {

	transport := http.DefaultTransport
	transport = roundtripper.NewRetryRoundTripper(transport, []int{420, 429}, 3)
	httpClient := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	productService := service.NewProductServiceHTTP(
		httpClient,
		app.config.ProductService.Token,
		fmt.Sprintf("%s://%s:%s", app.config.ProductService.Protocol, app.config.ProductService.Host, app.config.ProductService.Port),
	)

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", app.config.LomsService.Host, app.config.LomsService.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc.NewClient: %w", err)
	}
	stockClient := stocks.NewStockServiceClient(conn)
	orderClient := orders.NewOrderServiceClient(conn)
	lomsService := service.NewLomsServiceGRPC(stockClient, orderClient)

	const cartsStorageCap = 100
	cartRepository := repository.NewInMemoryCartRepository(cartsStorageCap)

	cartService := service.NewCartService(cartRepository, productService, lomsService)

	s := handler.NewServer(cartService, lomsService)

	mx := http.NewServeMux()
	mx.HandleFunc("POST /user/{user_id}/cart/{sku_id}", s.AddCartItemHandler)
	mx.HandleFunc("DELETE /user/{user_id}/cart/{sku_id}", s.DeleteCartItemHandler)
	mx.HandleFunc("DELETE /user/{user_id}/cart", s.ClearCartHandler)
	mx.HandleFunc("GET /user/{user_id}/cart", s.GetCartHandler)
	mx.HandleFunc("POST /checkout/{user_id}", s.CheckoutCartHandler)

	h := middleware.NewLoggerMiddleware(mx)

	return h, nil
}
