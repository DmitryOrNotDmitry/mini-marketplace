package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config главный конфиг сервиса.
type Config struct {
	Server         CartServiceConfig    `yaml:"service"`
	ProductService ProductServiceConfig `yaml:"product_service"`
	LomsService    LomsServiceConfig    `yaml:"loms_service"`
	Jaeger         JaegerConfig         `yaml:"jaeger"`
	RepoObserver   RepoObserverConfig   `yaml:"repo_observer"`
}

// CartServiceConfig конфиг для сервиса cart.
type CartServiceConfig struct {
	Host                    string        `yaml:"host"`
	Port                    string        `yaml:"port"`
	GracefulShutdownTimeout string        `yaml:"graceful_shutdown_timeout"`
	Tracing                 TracingConfig `yaml:"tracing"`
}

// ProductServiceConfig конфиг для сервиса product.
type ProductServiceConfig struct {
	Protocol string `yaml:"protocol"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Token    string `yaml:"token"`
	Limit    int    `yaml:"limit"`
}

// JaegerConfig конфиг для jaeger.
type JaegerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// TracingConfig конфиг для трассировки.
type TracingConfig struct {
	ServiceName string `yaml:"service_name"`
	Environment string `yaml:"environment"`
}

// LomsServiceConfig конфиг для сервиса loms.
type LomsServiceConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// RepoObserverConfig конфиг для трассировки.
type RepoObserverConfig struct {
	Interval int `yaml:"interval"`
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
