package flagsparser

import (
	fl "flag"

	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
)

type flagsParser struct {
}

// Parse парсит аргументы командной строки и наполняет конфигурацию
func (p flagsParser) Parse(c *config.Config) error {
	fl.StringVar(&c.Host, "a", c.Host, "Net address host:port")
	fl.Int64Var(&c.StoreInterval, "i", c.StoreInterval, "Store interval")
	fl.StringVar(&c.FileStoragePath, "f", c.FileStoragePath, "File storage path (full)")
	fl.BoolVar(&c.Restore, "r", c.Restore, "Restore values on start")
	fl.StringVar(&c.DatabaseDSN, "d", c.DatabaseDSN, "Database DSN")
	fl.StringVar(&c.Key, "k", c.Key, "Hash key")
	fl.StringVar(&c.CryptoKey, "crypto-key", c.CryptoKey, "Path to file with private key")
	fl.StringVar(&c.TrustedSubnet, "t", c.TrustedSubnet, "Строковое представление бесклассовой адресации (CIDR)")
	fl.Parse()
	return nil
}

// New возвращает сервис flagsParser для парсинга аргументов командной строки
func New() *flagsParser {
	return &flagsParser{}
}
