package configpathprovider

import (
	"github.com/caarlos0/env/v6"
)

type envConfig struct {
	Config string `env:"CONFIG"`
}

// FromEnvProvider сервис для получения пути к файлу конфигурации из параметров среды
type FromEnvProvider struct {
}

// Fetch метод получения пути к файлу конфигурации
func (p FromEnvProvider) Fetch() (string, error) {
	cfg := envConfig{Config: ""}
	err := env.Parse(&cfg)
	return cfg.Config, err
}
