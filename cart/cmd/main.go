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
	logLevel      = zap.InfoLevel
	serviceName   = "cart"
	configPathVar = "CONFIG_FILE"
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
		err := cartApp.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	pkgapp.GracefulShutdown(&shutdownMain{cancel: cancel, app: cartApp}, time.Duration(cartApp.Config.Server.GracefulShutdownTimeout)*time.Second)
}
