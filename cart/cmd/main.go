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

const (
	logLevel      = zap.InfoLevel
	serviceName   = "cart"
	configPathVar = "CONFIG_FILE"
)

func main() {
	logger.InitLogger(&logger.Config{
		Level:       logLevel,
		ServiceName: serviceName,
	})

	cartApp, err := app.NewApp(os.Getenv(configPathVar))
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
