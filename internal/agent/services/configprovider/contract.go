package configprovider

import "github.com/AndrXxX/go-metrics-collector/internal/agent/config"

type parser interface {
	Parse(c *config.Config) error
}
