package envparser

import (
	"fmt"

	"github.com/caarlos0/env/v6"

	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
)

type envConfig struct {
	Addr            string `env:"ADDRESS"`
	GRPCAddr        string `env:"GRPC_ADDRESS"`
	StoreInterval   int64  `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	Key             string `env:"KEY"`
	CryptoKey       string `env:"CRYPTO_KEY"`
	TrustedSubnet   string `env:"TRUSTED_SUBNET"`
}

type envParser struct {
}

// Parse парсит переменные окружения и наполняет конфигурацию
func (p envParser) Parse(c *config.Config) error {
	cfg := envConfig{
		Addr:            c.Host,
		GRPCAddr:        c.GRPCHost,
		StoreInterval:   c.StoreInterval,
		FileStoragePath: c.FileStoragePath,
		Restore:         c.Restore,
		DatabaseDSN:     c.DatabaseDSN,
		Key:             c.Key,
		CryptoKey:       c.CryptoKey,
		TrustedSubnet:   c.TrustedSubnet,
	}
	err := env.Parse(&cfg)
	if err != nil {
		return fmt.Errorf("error on parse config: %w", err)
	}
	c.Host = cfg.Addr
	c.GRPCHost = cfg.GRPCAddr
	c.StoreInterval = cfg.StoreInterval
	c.FileStoragePath = cfg.FileStoragePath
	c.Restore = cfg.Restore
	c.DatabaseDSN = cfg.DatabaseDSN
	c.Key = cfg.Key
	c.CryptoKey = cfg.CryptoKey
	c.TrustedSubnet = cfg.TrustedSubnet
	return nil
}

// New возвращает сервис envParser для парсинга переменных окружения
func New() *envParser {
	return &envParser{}
}
