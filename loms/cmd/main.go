package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"route256/cart/pkg/myerrgroup"
	"route256/loms/internal/app"
	pkgapp "route256/loms/pkg/app"
)

func main() {
	lomsApp, err := app.NewApp(os.Getenv("CONFIG_FILE"))
	if err != nil {
		panic(err)
	}

	go func() {
		err = startApp(lomsApp)
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(lomsApp.Config.Server.GracefullShutdownTimeout)*time.Second)
	defer cancel()

	pkgapp.GracefullShutdown(ctx, lomsApp)
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
