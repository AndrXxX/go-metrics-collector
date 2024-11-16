package configfile

import (
	"github.com/AndrXxX/go-metrics-collector/internal/types/jsontime"
)

type jsonConfig struct {
	Address       *string            `json:"address"`        // аналог переменной окружения ADDRESS или флага -a
	Restore       *bool              `json:"restore"`        // аналог переменной окружения RESTORE или флага -r
	StoreInterval *jsontime.Duration `json:"store_interval"` // аналог переменной окружения STORE_INTERVAL или флага -i
	StoreFile     *string            `json:"store_file"`     // аналог переменной окружения STORE_FILE или -f
	DatabaseDsn   *string            `json:"database_dsn"`   // аналог переменной окружения DATABASE_DSN или флага -d
	CryptoKey     *string            `json:"crypto_key"`     // аналог переменной окружения CRYPTO_KEY или флага -crypto-key
}
