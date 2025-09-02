package main

import (
	"os"
	"route256/cart/internal/app"
)

func main() {
	cartApp, err := app.NewApp(os.Getenv("CONFIG_FILE"))
	if err != nil {
		panic(err)
	}

	if err := cartApp.ListenAndServe(); err != nil {
		panic(err)
	}
}
