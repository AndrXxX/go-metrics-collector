package envparser

import (
	"strings"

	"github.com/caarlos0/env/v6"

	"github.com/AndrXxX/go-metrics-collector/internal/staticlint/config"
)

type envConfig struct {
	Analyzers string `env:"STATIC_ANALYZERS"`
}

// EnvParser сервис для парсинга переменных окружения staticlint
type EnvParser struct {
}

// Parse парсит переменные окружения и наполняет конфигурацию для staticlint
func (p EnvParser) Parse(c *config.Config) error {
	cfg := envConfig{
		Analyzers: strings.Join(c.StaticAnalyzers, ","),
	}
	err := env.Parse(&cfg)
	if err != nil {
		return err
	}
	c.StaticAnalyzers = strings.Split(cfg.Analyzers, ",")
	return nil
}
