package main

import (
	"os"
	"route256/cart/internal/app"
)

func main() {
	cfg := os.Getenv("ROUTE_256_WS1_CONFIG")
	if cfg == "" {
		cfg = "cart/configs/values_ci.yaml"
	}
	app, err := app.NewApp(cfg)
	if err != nil {
		panic(err)
	}

	if err := app.ListenAndServe(); err != nil {
		panic(err)
	}
}
