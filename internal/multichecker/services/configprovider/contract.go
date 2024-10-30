package configprovider

import "github.com/AndrXxX/go-metrics-collector/internal/multichecker/config"

type parser interface {
	Parse(c *config.Config) error
}
