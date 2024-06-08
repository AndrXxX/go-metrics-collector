package main

import (
	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
	"github.com/caarlos0/env/v6"
)

type envConfig struct {
	addr string `env:"ADDRESS"`
}

func parseEnv(c *config.Config) {
	var cfg envConfig
	err := env.Parse(&cfg)
	if err != nil {
		return
	}
	if cfg.addr != "" {
		c.Host = cfg.addr
	}
}
