package configfile

import (
	"github.com/AndrXxX/go-metrics-collector/internal/types/jsontime"
)

type jsonConfig struct {
	Address        *string            `json:"address"`         // Аналог переменной окружения ADDRESS или флага -a
	ReportInterval *jsontime.Duration `json:"report_interval"` // Аналог переменной окружения REPORT_INTERVAL или флага -r
	PollInterval   *jsontime.Duration `json:"poll_interval"`   // Аналог переменной окружения POLL_INTERVAL или флага -p
	CryptoKey      *string            `json:"crypto_key"`      // Аналог переменной окружения CRYPTO_KEY или флага -crypto-key
}
