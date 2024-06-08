package main

import (
	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
	"github.com/caarlos0/env/v6"
)

type envConfig struct {
	addr           string `env:"ADDRESS"`
	reportInterval int64  `env:"REPORT_INTERVAL"`
	pollInterval   int64  `env:"POLL_INTERVAL"`
}

func parseEnv(c *config.Config) {
	cfg := envConfig{
		addr:           c.Common.Host,
		reportInterval: c.Intervals.ReportInterval,
		pollInterval:   c.Intervals.PollInterval,
	}
	err := env.Parse(&cfg)
	if err != nil {
		return
	}
	if cfg.addr != "" {
		c.Common.Host = cfg.addr
	}
	c.Intervals.ReportInterval = cfg.reportInterval
	c.Intervals.PollInterval = cfg.pollInterval
}
