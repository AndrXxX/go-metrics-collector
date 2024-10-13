package main

import (
	"github.com/caarlos0/env/v6"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
)

type EnvConfig struct {
	Addr           string `env:"ADDRESS"`
	ReportInterval int64  `env:"REPORT_INTERVAL"`
	PollInterval   int64  `env:"POLL_INTERVAL"`
	Key            string `env:"KEY"`
	RateLimit      int64  `env:"RATE_LIMIT"`
}

func parseEnv(c *config.Config) error {
	cfg := EnvConfig{
		Addr:           c.Common.Host,
		Key:            c.Common.Key,
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
	c.Common.RateLimit = cfg.RateLimit
	c.Intervals.ReportInterval = cfg.ReportInterval
	c.Intervals.PollInterval = cfg.PollInterval
	return nil
}
