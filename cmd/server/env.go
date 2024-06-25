package main

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/logger"
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
		logger.Log.Error(fmt.Sprintf("Error on parse EnvConfig: %s", err.Error()))
		return
	}
	c.Host = cfg.Addr
}
