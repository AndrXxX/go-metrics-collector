package main

import (
	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
	"github.com/caarlos0/env/v6"
)

type EnvConfig struct {
	Addr string `env:"ADDRESS"`
}

func parseEnv(c *config.Config) {
	cfg := EnvConfig{
		Addr: c.Host,
	}
	err := env.Parse(&cfg)
	if err != nil {
		return
	}
	c.Host = cfg.Addr
}
