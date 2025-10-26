package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config главный конфиг сервиса.
type Config struct {
	Server CommentsServiceConfig `yaml:"service"`
	App    AppConfig             `yaml:"app"`
	DB     DBConfig              `yaml:"db"`
}

// CommentsServiceConfig конфиг для сервиса comments.
type CommentsServiceConfig struct {
	Host                    string            `yaml:"host"`
	GRPCPort                string            `yaml:"grpc_port"`
	HTTPPort                string            `yaml:"http_port"`
	GRPCGateWay             GRPCGateWayConfig `yaml:"grpc_gateway"`
	GracefulShutdownTimeout string            `yaml:"graceful_shutdown_timeout"`
}

// AppConfig конфиг для настроек приложения.
type AppConfig struct {
	EditInterval   string `yaml:"edit_interval"`
	LimitRowsBySku int32  `yaml:"limit_rows_by_sku"`
}

// GRPCGateWayConfig конфиг для gRPC-gateway.
type GRPCGateWayConfig struct {
	ReadHeaderTimeout string `yaml:"read_header_timeout"`
	WriteTimeout      string `yaml:"write_timeout"`
	IdleTimeout       string `yaml:"idle_timeout"`
}

// DBConfig конфиг для БД
type DBConfig struct {
	Buckets int64           `yaml:"buckets"`
	Shards  []DBShardConfig `yaml:"shards"`
}

// DBShardConfig конфиг для шарда БД
type DBShardConfig struct {
	Host           string `yaml:"host"`
	Port           int64  `yaml:"port"`
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	DBName         string `yaml:"db_name"`
	BucketPosition int64  `yaml:"bucket_position"`
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
