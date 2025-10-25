package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config главный конфиг сервиса.
type Config struct {
	Server    LomsServiceConfig `yaml:"service"`
	MasterDB  MasterDBConfig    `yaml:"db_master"`
	ReplicaDB ReplicaDBConfig   `yaml:"db_replica"`
	Jaeger    JaegerConfig      `yaml:"jaeger"`
}

// LomsServiceConfig конфиг для сервиса loms.
type LomsServiceConfig struct {
	Host                    string            `yaml:"host"`
	GRPCPort                string            `yaml:"grpc_port"`
	HTTPPort                string            `yaml:"http_port"`
	GRPCGateWay             GRPCGateWayConfig `yaml:"grpc_gateway"`
	GracefulShutdownTimeout string            `yaml:"graceful_shutdown_timeout"`
	Tracing                 TracingConfig     `yaml:"tracing"`
}

// GRPCGateWayConfig конфиг для gRPC-gateway.
type GRPCGateWayConfig struct {
	ReadHeaderTimeout string `yaml:"read_header_timeout"`
	WriteTimeout      string `yaml:"write_timeout"`
	IdleTimeout       string `yaml:"idle_timeout"`
}

// MasterDBConfig конфиг для мастера БД
type MasterDBConfig struct {
	Host     string `yaml:"host"`
	Port     int64  `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
}

// ReplicaDBConfig конфиг для реплики БД.
type ReplicaDBConfig struct {
	Host     string `yaml:"host"`
	Port     int64  `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
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
