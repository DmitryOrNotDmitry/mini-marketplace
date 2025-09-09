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

	err = lomsApp.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
