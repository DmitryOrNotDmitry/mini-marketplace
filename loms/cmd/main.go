package main

import (
	"os"
	"route256/loms/internal/app"

	"golang.org/x/sync/errgroup"
)

func main() {
	lomsApp, err := app.NewApp(
		os.Getenv("CONFIG_FILE"),
	)
	if err != nil {
		panic(err)
	}

	errGroup := new(errgroup.Group)
	errGroup.Go(func() error {
		return lomsApp.ListenAndServeGRPCGateway()
	})

	errGroup.Go(func() error {
		return lomsApp.ListenAndServeGRPC()
	})

	if err := errGroup.Wait(); err != nil {
		panic(err)
	}
}
