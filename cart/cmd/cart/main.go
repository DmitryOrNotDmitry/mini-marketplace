package main

import (
	"os"
	"route256/cart/internal/app"
)

func main() {
	app, err := app.NewApp(os.Getenv("ROUTE_256_WS1_CONFIG"))
	if err != nil {
		panic(err)
	}

	if err := app.ListenAndServe(); err != nil {
		panic(err)
	}
}
