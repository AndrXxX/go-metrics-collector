package envparser

import (
	"github.com/caarlos0/env/v6"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
)

type envConfig struct {
	Addr           string `env:"ADDRESS"`
	ReportInterval int64  `env:"REPORT_INTERVAL"`
	PollInterval   int64  `env:"POLL_INTERVAL"`
	Key            string `env:"KEY"`
	CryptoKey      string `env:"CRYPTO_KEY"`
	RateLimit      int64  `env:"RATE_LIMIT"`
}

type envParser struct {
}

// Parse парсит переменные окружения и наполняет конфигурацию
func (p envParser) Parse(c *config.Config) error {
	cfg := envConfig{
		Addr:           c.Common.Host,
		Key:            c.Common.Key,
		CryptoKey:      c.Common.CryptoKey,
		RateLimit:      c.Common.RateLimit,
		ReportInterval: c.Intervals.ReportInterval,
		PollInterval:   c.Intervals.PollInterval,
	}
	err := env.Parse(&cfg)
	if err != nil {
		return err
	}
	c.Common.Host = cfg.Addr
	c.Common.Key = cfg.Key
	c.Common.CryptoKey = cfg.CryptoKey
	c.Common.RateLimit = cfg.RateLimit
	c.Intervals.ReportInterval = cfg.ReportInterval
	c.Intervals.PollInterval = cfg.PollInterval
	return nil
}

// New возвращает сервис envParser для парсинга переменных окружения
func New() *envParser {
	return &envParser{}
}
