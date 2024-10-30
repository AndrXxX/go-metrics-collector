package configprovider

import "github.com/AndrXxX/go-metrics-collector/internal/staticlint/config"

type parser interface {
	Parse(c *config.Config) error
}
