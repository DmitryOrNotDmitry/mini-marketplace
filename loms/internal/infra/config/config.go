package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server LomsServiceConfig `yaml:"service"`
}

type LomsServiceConfig struct {
	Host     string `yaml:"host"`
	GRPCPort string `yaml:"grpc_port"`
	HTTPPort string `yaml:"http_port"`
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
