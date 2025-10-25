package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"route256/cart/pkg/logger"
	"route256/cart/pkg/myerrgroup"
	"route256/comments/internal/app"
	pkgapp "route256/loms/pkg/app"

	"go.uber.org/zap"
)

const (
	logLevel                       = zap.InfoLevel
	serviceName                    = "comments"
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
	commentsApp, err := app.NewApp(mainCtx, os.Getenv(configPathVar))
	if err != nil {
		panic(err)
	}

	go func() {
		err = startApp(mainCtx, commentsApp)
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	gracefulTimeout, err := time.ParseDuration(commentsApp.Config.Server.GracefulShutdownTimeout)
	if err != nil {
		gracefulTimeout = defaultGracefulShutdownTimeout
	}
	pkgapp.GracefulShutdown(&shutdownMain{cancel: cancel, app: commentsApp}, gracefulTimeout)
}

func startApp(ctx context.Context, commentsApp *app.App) error {
	errGroup := myerrgroup.New()
	errGroup.Go(func() error {
		return commentsApp.ListenAndServeGRPCGateway(ctx)
	})

	errGroup.Go(func() error {
		return commentsApp.ListenAndServeGRPC()
	})

	return errGroup.Wait()
}
