package config

import "github.com/AndrXxX/go-metrics-collector/internal/enums/vars"

func NewConfig() *Config {
	return &Config{
		Host:     vars.DefaultHost,
		LogLevel: vars.DefaultLogLevel,
	}
}
