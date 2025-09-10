package main

import (
	"os"
	"route256/loms/internal/app"
)

func main() {
	lomsApp, err := app.NewApp(os.Getenv("CONFIG_FILE"))
	if err != nil {
		panic(err)
	}

	go func() {
		err = lomsApp.ListenAndServeGRPCGateway()
		if err != nil {
			panic(err)
		}
	}()

	err = lomsApp.ListenAndServeGRPC()
	if err != nil {
		panic(err)
	}
}
