package config

import "github.com/AndrXxX/go-metrics-collector/internal/multichecker/vars"

// NewConfig предоставляет конфигурацию со значениями по умолчанию для multichecker
func NewConfig() *Config {
	return &Config{
		StaticChecks: []string{
			vars.StaticSAChecks,
			vars.StaticSTChecks,
			vars.StaticQFChecks,
		},
	}
}
