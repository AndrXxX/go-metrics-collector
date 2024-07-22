package config

import "github.com/AndrXxX/go-metrics-collector/internal/enums/vars/defaults"

func NewConfig() *Config {
	return &Config{
		Host:            defaults.Host,
		LogLevel:        defaults.LogLevel,
		StoreInterval:   defaults.StoreInterval,
		RepeatIntervals: defaults.RepeatIntervals,
		FileStoragePath: defaults.FileStoragePath,
		Restore:         defaults.Restore,
	}
}
