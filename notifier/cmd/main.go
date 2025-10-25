package main

import (
	"context"
	"os"
	"time"

	"route256/cart/pkg/logger"
	pkgapp "route256/loms/pkg/app"
	"route256/notifier/internal/app"

	"go.uber.org/zap"
)

const (
	logLevel      = zap.InfoLevel
	serviceName   = "notifier"
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
	notifierApp, err := app.NewApp(mainCtx, os.Getenv(configPathVar))
	if err != nil {
		panic(err)
	}

	pkgapp.GracefulShutdown(&shutdownMain{cancel: cancel, app: notifierApp}, time.Duration(notifierApp.Config.Server.GracefulShutdownTimeout)*time.Second)
}
