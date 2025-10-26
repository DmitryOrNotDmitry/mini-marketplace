package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config главный конфиг сервиса.
type Config struct {
	Server ServiceConfig `yaml:"service"`
	Kafka  KafkaConfig   `yaml:"kafka"`
}

// KafkaConfig конфиг для kafka.
type KafkaConfig struct {
	OrderTopic      string `yaml:"order_topic"`
	ConsumerGroupID string `yaml:"consumer_group_id"`
	Brokers         string `yaml:"brokers"`
}

// Config конфиг для текущего сервиса.
type ServiceConfig struct {
	GracefulShutdownTimeout int64 `yaml:"graceful_shutdown_timeout"`
}

// LoadConfig загружает конфиг из файла .yaml
func LoadConfig(filename string) (*Config, error) {
	f, err := os.Open(filename) // nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("ошибка при открытии файла-конфига: %w", err)
	}
	defer f.Close()

	config := &Config{}
	if err := yaml.NewDecoder(f).Decode(config); err != nil {
		return nil, fmt.Errorf("ошибка при декодировании yaml файла-конфига: %w", err)
	}

	return config, nil
}
