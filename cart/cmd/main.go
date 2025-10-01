package main

import (
	"net/http"
	"os"
	"time"

	"route256/cart/internal/app"
	"route256/cart/pkg/logger"
	pkgapp "route256/loms/pkg/app"

	"go.uber.org/zap"
)

func main() {
	logger.InitLogger(&logger.Config{
		Level:       zap.InfoLevel,
		ServiceName: "cart",
	})

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
