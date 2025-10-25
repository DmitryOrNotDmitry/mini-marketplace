package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"route256/cart/internal/app"
	"route256/cart/pkg/logger"
	pkgapp "route256/loms/pkg/app"

	"go.uber.org/zap"
)

const (
	logLevel                       = zap.InfoLevel
	serviceName                    = "cart"
	configPathVar                  = "CONFIG_FILE"
	defaultGracefulShutdownTimeout = 10 * time.Second
)

type shutdownMain struct {
	cancel context.CancelFunc
	app    *app.App
}

func (a *shutdownMain) Shutdown(ctx context.Context) error {
	a.cancel()
	return a.app.Shutdown(ctx)
}

func main() {
	logger.InitLogger(&logger.Config{
		Level:       logLevel,
		ServiceName: serviceName,
	})

	mainCtx, cancel := context.WithCancel(context.Background())
	cartApp, err := app.NewApp(mainCtx, os.Getenv(configPathVar))
	if err != nil {
		panic(err)
	}

	go func() {
		httpServerErr := cartApp.ListenAndServe()
		if httpServerErr != nil && httpServerErr != http.ErrServerClosed {
			panic(httpServerErr)
		}
	}()

	gracefulTimeout, err := time.ParseDuration(cartApp.Config.Server.GracefulShutdownTimeout)
	if err != nil {
		gracefulTimeout = defaultGracefulShutdownTimeout
	}
	pkgapp.GracefulShutdown(&shutdownMain{cancel: cancel, app: cartApp}, gracefulTimeout)
}
