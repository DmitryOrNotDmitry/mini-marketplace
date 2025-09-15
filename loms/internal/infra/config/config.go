package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config главный конфиг сервиса.
type Config struct {
	Server LomsServiceConfig `yaml:"service"`
}

// LomsServiceConfig конфиг для сервиса loms.
type LomsServiceConfig struct {
	Host                     string            `yaml:"host"`
	GRPCPort                 string            `yaml:"grpc_port"`
	HTTPPort                 string            `yaml:"http_port"`
	GRPCGateWay              GRPCGateWayConfig `yaml:"grpc_gateway"`
	LoadStocksDataTimeout    int64             `yaml:"load_stocks_data_timeout"`
	GracefullShutdownTimeout int64             `yaml:"gracefull_shutdown_timeout"`
}

// GRPCGateWayConfig конфиг для gRPC-gateway.
type GRPCGateWayConfig struct {
	ReadHeaderTimeout int64 `yaml:"read_header_timeout"`
	WriteTimeout      int64 `yaml:"write_timeout"`
	IdleTimeout       int64 `yaml:"idle_timeout"`
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
