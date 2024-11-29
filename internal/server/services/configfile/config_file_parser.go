package configfile

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
	"github.com/AndrXxX/go-metrics-collector/internal/types/jsontime"
)

// Parser сервис для парсинга файла конфигурации
type Parser struct {
	PathProvider pathProvider
}

// Parse парсит файл конфигурации
func (p Parser) Parse(c *config.Config) error {
	if p.PathProvider == nil {
		return nil
	}
	path, err := p.PathProvider.Fetch()
	if err != nil {
		return err
	}
	if path == "" {
		return nil
	}
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open config file: %s", err)
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	jc := jsonConfig{}
	err = json.NewDecoder(f).Decode(&jc)
	if err != nil {
		return fmt.Errorf("failed to parse config file: %s", err)
	}
	set(jc.Address, &c.Host)
	set(jc.Restore, &c.Restore)
	set(convertDurationToInt(jc.StoreInterval), &c.StoreInterval)
	set(jc.StoreFile, &c.FileStoragePath)
	set(jc.DatabaseDsn, &c.DatabaseDSN)
	set(jc.CryptoKey, &c.CryptoKey)
	set(jc.TrustedSubnet, &c.TrustedSubnet)
	return nil
}

func set[T comparable](val *T, target *T) {
	var zero T
	if val == nil || *val == zero {
		return
	}
	*target = *val
}

func convertDurationToInt(d *jsontime.Duration) *int64 {
	var zero jsontime.Duration
	if d == nil || *d == zero {
		return nil
	}
	v := int64(d.Duration / time.Second)
	return &v
}
