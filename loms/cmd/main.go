package main

import (
	"net/http"
	"os"
	"time"

	"route256/cart/pkg/logger"
	"route256/cart/pkg/myerrgroup"
	"route256/loms/internal/app"
	pkgapp "route256/loms/pkg/app"

	"go.uber.org/zap"
)

const (
	logLevel      = zap.InfoLevel
	serviceName   = "loms"
	configPathVar = "CONFIG_FILE"
)

func main() {
	logger.InitLogger(&logger.Config{
		Level:       logLevel,
		ServiceName: serviceName,
	})

	lomsApp, err := app.NewApp(os.Getenv(configPathVar))
	if err != nil {
		panic(err)
	}

	go func() {
		err = startApp(lomsApp)
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	pkgapp.GracefulShutdown(lomsApp, time.Duration(lomsApp.Config.Server.GracefulShutdownTimeout)*time.Second)
}

func startApp(lomsApp *app.App) error {
	errGroup := myerrgroup.New()
	errGroup.Go(func() error {
		return lomsApp.ListenAndServeGRPCGateway()
	})

	errGroup.Go(func() error {
		return lomsApp.ListenAndServeGRPC()
	})

	return errGroup.Wait()
}
