package configprovider

import "github.com/AndrXxX/go-metrics-collector/internal/server/config"

type flagsParser interface {
	Parse(c *config.Config)
}

type envParser interface {
	Parse(c *config.Config) error
}
