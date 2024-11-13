package config

// Config конфигурация сервера
type Config struct {
	Host            string `valid:"minstringlength(3)"`
	LogLevel        string `valid:"in(debug|info|warn|error|fatal)"`
	StoreInterval   int    `valid:"range(1|999)"`
	RepeatIntervals []int
	FileStoragePath string `valid:"minstringlength(3)"`
	Restore         bool
	DatabaseDSN     string
	Key             string
	CryptoKey       string
}
