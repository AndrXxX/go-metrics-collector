package config

import "github.com/AndrXxX/go-metrics-collector/internal/enums/vars/defaults"

func NewConfig() *Config {
	return &Config{
		Host:            defaults.Host,
		LogLevel:        defaults.LogLevel,
		StoreInterval:   defaults.StoreInterval,
		FileStoragePath: defaults.FileStoragePath,
	}
}
