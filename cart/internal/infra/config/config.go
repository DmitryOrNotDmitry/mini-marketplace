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
}

// CartServiceConfig конфиг для сервиса cart.
type CartServiceConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// ProductServiceConfig конфиг для сервиса product.
type ProductServiceConfig struct {
	Protocol string `yaml:"protocol"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Token    string `yaml:"token"`
}

// LomsServiceConfig конфиг для сервиса loms.
type LomsServiceConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
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
