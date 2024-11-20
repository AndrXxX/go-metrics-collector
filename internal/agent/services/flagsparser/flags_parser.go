package flagsparser

import (
	fl "flag"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
)

type flagsParser struct {
}

// Parse парсит аргументы командной строки и наполняет конфигурацию
func (p flagsParser) Parse(c *config.Config) error {
	fl.StringVar(&c.Common.Host, "a", c.Common.Host, "Net address host:port")
	fl.StringVar(&c.Common.Key, "k", c.Common.Key, "Hash key")
	fl.StringVar(&c.Common.CryptoKey, "crypto-key", c.Common.CryptoKey, "Path to file with public key")
	fl.Int64Var(&c.Common.RateLimit, "l", c.Common.RateLimit, "Rate Limit")
	fl.Int64Var(&c.Intervals.ReportInterval, "r", c.Intervals.ReportInterval, "Report Interval in seconds")
	fl.Int64Var(&c.Intervals.PollInterval, "p", c.Intervals.PollInterval, "Poll Interval in seconds")
	fl.Parse()
	return nil
}

// New возвращает сервис flagsParser для парсинга аргументов командной строки
func New() *flagsParser {
	return &flagsParser{}
}
