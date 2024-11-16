package configfile

import (
	"encoding/json"
	"fmt"
	"time"
)

type jsonConfig struct {
	Address        *string   `json:"address"`         // Аналог переменной окружения ADDRESS или флага -a
	ReportInterval *duration `json:"report_interval"` // Аналог переменной окружения REPORT_INTERVAL или флага -r
	PollInterval   *duration `json:"poll_interval"`   // Аналог переменной окружения POLL_INTERVAL или флага -p
	CryptoKey      *string   `json:"crypto_key"`      // Аналог переменной окружения CRYPTO_KEY или флага -crypto-key
}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case string:
		dVal, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = duration{dVal}
		return nil
	}
	return fmt.Errorf("cannot unmarshal JSON value")
}
