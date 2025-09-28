package main

import (
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

	pkgapp.GracefulShutdown(cartApp, time.Duration(cartApp.Config.Server.GracefulShutdownTimeout)*time.Second)
}
