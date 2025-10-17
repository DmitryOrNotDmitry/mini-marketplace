package app

import (
	"context"
	"fmt"
	"route256/notifier/internal/infra/config"
	"route256/notifier/internal/infra/kafka"
	"route256/notifier/internal/service"
)

// App создает компоненты для сервиса notifier
type App struct {
	Config *config.Config
}

// NewApp конструктор главного приложения.
func NewApp(ctx context.Context, configPath string) (*App, error) {
	c, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("config.LoadConfig: %w", err)
	}

	a := &App{Config: c}

	orderEventConsumer, err := kafka.NewOrderEventTopicSubKafka(a.Config.Kafka.ConsumerGroupID, []string{a.Config.Kafka.OrderTopic}, []string{a.Config.Kafka.Brokers})
	if err != nil {
		return nil, fmt.Errorf("config.LoadConfig: %w", err)
	}

	msgIn := orderEventConsumer.Start(ctx)

	orderEventService := service.NewOrderEventService(msgIn)
	orderEventService.Start()

	return a, nil
}

// Shutdown gracefully останавливает приложение.
func (a *App) Shutdown(_ context.Context) error {
	return nil
}
