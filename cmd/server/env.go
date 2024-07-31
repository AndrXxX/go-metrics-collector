package main

import (
	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
)

type EnvConfig struct {
	Addr            string `env:"ADDRESS"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	Key             string `env:"KEY"`
}

func parseEnv(c *config.Config) {
	cfg := EnvConfig{
		Addr:            c.Host,
		StoreInterval:   c.StoreInterval,
		FileStoragePath: c.FileStoragePath,
		Restore:         c.Restore,
		DatabaseDSN:     c.DatabaseDSN,
		Key:             c.Key,
	}
	err := env.Parse(&cfg)
	if err != nil {
		logger.Log.Error("Error on parse EnvConfig", zap.Error(err))
		return
	}
	c.Host = cfg.Addr
	c.StoreInterval = cfg.StoreInterval
	c.FileStoragePath = cfg.FileStoragePath
	c.Restore = cfg.Restore
	c.DatabaseDSN = cfg.DatabaseDSN
	c.Key = cfg.Key
}
