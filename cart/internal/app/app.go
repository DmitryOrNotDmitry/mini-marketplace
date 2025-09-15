package app

import (
	"context"
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

// App создает компоненты для сервиса cart
type App struct {
	Config *config.Config
	server http.Server
}

// NewApp конструктор главного приложения.
func NewApp(configPath string) (*App, error) {
	c, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("config.LoadConfig: %w", err)
	}

	app := &App{Config: c}
	app.server.Handler, err = app.bootstrapHandlers()
	if err != nil {
		return nil, fmt.Errorf("app.bootstrapHandlers: %w", err)
	}

	return app, nil
}

// ListenAndServe запускает HTTP-сервер приложения.
func (app *App) ListenAndServe() error {
	address := fmt.Sprintf("%s:%s", app.Config.Server.Host, app.Config.Server.Port)

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
		app.Config.ProductService.Token,
		fmt.Sprintf("%s://%s:%s", app.Config.ProductService.Protocol, app.Config.ProductService.Host, app.Config.ProductService.Port),
	)

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", app.Config.LomsService.Host, app.Config.LomsService.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc.NewClient: %w", err)
	}
	stockClient := stocks.NewStockServiceV1Client(conn)
	orderClient := orders.NewOrderServiceV1Client(conn)
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

// Shutdown gracefully останавливает приложение.
func (app *App) Shutdown(ctx context.Context) error {
	return app.server.Shutdown(ctx)
}
