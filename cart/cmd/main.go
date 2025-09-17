package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"route256/cart/internal/app"
	pkgapp "route256/loms/pkg/app"
)

func main() {
	cartApp, err := app.NewApp(os.Getenv("CONFIG_FILE"))
	if err != nil {
		panic(err)
	}

	go func() {
		err := cartApp.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cartApp.Config.Server.GracefullShutdownTimeout)*time.Second)
	defer cancel()

	pkgapp.GracefullShutdown(ctx, cartApp)
}
