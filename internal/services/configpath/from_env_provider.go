package configpath

import (
	"github.com/caarlos0/env/v6"
)

type envConfig struct {
	Config string `env:"CONFIG"`
}

type fromEnvProvider struct {
}

// Fetch метод получения пути к файлу конфигурации
func (p fromEnvProvider) Fetch() (string, error) {
	cfg := envConfig{Config: ""}
	err := env.Parse(&cfg)
	return cfg.Config, err
}
