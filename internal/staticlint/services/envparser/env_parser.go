package envparser

import (
	"strings"

	"github.com/caarlos0/env/v6"

	"github.com/AndrXxX/go-metrics-collector/internal/staticlint/config"
)

type envConfig struct {
	Analyzers        string `env:"STATIC_ANALYZERS"`
	ExcludeAnalyzers string `env:"EXCLUDE_STATIC_ANALYZERS"`
}

// EnvParser сервис для парсинга переменных окружения staticlint
type EnvParser struct {
}

// Parse парсит переменные окружения и наполняет конфигурацию для staticlint
func (p EnvParser) Parse(c *config.Config) error {
	cfg := envConfig{
		Analyzers:        strings.Join(c.StaticAnalyzers, ","),
		ExcludeAnalyzers: strings.Join(c.ExcludeStaticAnalyzers, ","),
	}
	err := env.Parse(&cfg)
	if err != nil {
		return err
	}
	if cfg.Analyzers != "" {
		c.StaticAnalyzers = strings.Split(cfg.Analyzers, ",")
	}
	if cfg.ExcludeAnalyzers != "" {
		c.ExcludeStaticAnalyzers = strings.Split(cfg.ExcludeAnalyzers, ",")
	}
	return nil
}
