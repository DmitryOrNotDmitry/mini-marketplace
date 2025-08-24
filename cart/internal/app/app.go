package app

import (
	"fmt"
	"net"
	"net/http"

	"route256/cart/internal/infra/config"
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

	fmt.Println("app bootstrap")

	return app.server.Serve(l)
}

func (app *App) bootstrapHandlers() http.Handler {

	//transport := http.DefaultTransport
	//transport = round_trippers.NewLogRoundTripper(transport)
	// httpClient := http.Client{
	// 	Transport: transport,
	// 	Timeout:   10 * time.Second,
	// }

	// productService := product_service.NewProductService(
	// 	httpClient,
	// 	app.config.ProductService.Token,
	// 	fmt.Sprintf("%s:%s", app.config.ProductService.Host, app.config.ProductService.Port),
	// )

	// var counter atomic.Uint64
	// const reviewsCap = 100
	// reviewRepository := repository.NewInMemoryRepository(reviewsCap, &counter)
	// reviewService := service.NewReviewService(reviewRepository, productService)

	//s := handler.NewServer()

	mx := http.NewServeMux()
	// mx.HandleFunc("POST /products/{sku}/reviews", s.AddReviewHandler)
	// mx.HandleFunc("GET /products/{sku}/reviews", s.GetReviewsBySkuHandler)

	//h := middlewares.NewTimerMiddleware(mx)

	return mx
}
