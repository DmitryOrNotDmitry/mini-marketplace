package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"route256/cart/pkg/logger"
	"syscall"
)

// ShutdownedApp определяет методы для приложения, которое можно остановить gracefully
type ShutdownedApp interface {
	Shutdown(context.Context) error
}

// GracefullShutdown выполняет gracefull shutdown для приложения.
func GracefullShutdown(ctx context.Context, shutApp ShutdownedApp) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	logger.Info("Shutting down service...")

	err := shutApp.Shutdown(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("Service stopped with force %s", err.Error()))
	} else {
		logger.Info("Service stopped gracefully")
	}
}
