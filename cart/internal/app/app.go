package app

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"route256/cart/internal/handler"
	"route256/cart/internal/infra/config"
	"route256/cart/internal/infra/http/middleware"
	"route256/cart/internal/infra/http/round_tripper"
	"route256/cart/internal/infra/logger"
	"route256/cart/internal/infra/repository"
	"route256/cart/internal/service"
)

type App struct {
	config *config.Config
	server http.Server
}

func NewApp(configPath string) (*App, error) {
	c, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("config.LoadConfig: %w", err)
	}

	app := &App{config: c}
	app.server.Handler = app.bootstrapHandlers()

	return app, nil
}

func (app *App) ListenAndServe() error {
	address := fmt.Sprintf("%s:%s", app.config.Server.Host, app.config.Server.Port)

	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("Cart service listening at http://%s", address))

	return app.server.Serve(l)
}

func (app *App) bootstrapHandlers() http.Handler {

	transport := http.DefaultTransport
	transport = round_tripper.NewRetryRoundTripper(transport, []int{420, 429}, 3)
	httpClient := http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	productService := service.NewProductService(
		httpClient,
		app.config.ProductService.Token,
		fmt.Sprintf("%s://%s:%s", app.config.ProductService.Protocol, app.config.ProductService.Host, app.config.ProductService.Port),
	)

	const cartsStorageCap = 100
	cartRepository := repository.NewInMemoryCartRepository(cartsStorageCap)
	cartService := service.NewCartService(cartRepository, productService)

	s := handler.NewServer(cartService)

	mx := http.NewServeMux()
	mx.HandleFunc("POST /user/{user_id}/cart/{sku_id}", s.AddCartItemHandler)
	mx.HandleFunc("DELETE /user/{user_id}/cart/{sku_id}", s.DeleteCartItemHandler)
	mx.HandleFunc("DELETE /user/{user_id}/cart", s.ClearCartHandler)
	mx.HandleFunc("GET /user/{user_id}/cart", s.GetCartHandler)

	h := middleware.NewLoggerMiddleware(mx)

	return h
}
