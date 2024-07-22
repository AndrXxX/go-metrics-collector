package main

import (
	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
	"github.com/caarlos0/env/v6"
)

type EnvConfig struct {
	Addr           string `env:"ADDRESS"`
	ReportInterval int64  `env:"REPORT_INTERVAL"`
	PollInterval   int64  `env:"POLL_INTERVAL"`
	Key            string `env:"KEY"`
}

func parseEnv(c *config.Config) error {
	cfg := EnvConfig{
		Addr:           c.Common.Host,
		ReportInterval: c.Intervals.ReportInterval,
		PollInterval:   c.Intervals.PollInterval,
	}
	err := env.Parse(&cfg)
	if err != nil {
		return err
	}
	c.Common.Host = cfg.Addr
	c.Common.Key = cfg.Key
	c.Intervals.ReportInterval = cfg.ReportInterval
	c.Intervals.PollInterval = cfg.PollInterval
	return nil
}
