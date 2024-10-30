package configprovider

import (
	"fmt"

	"github.com/AndrXxX/go-metrics-collector/internal/multichecker/config"
)

type configProvider struct {
	parsers []parser
}

// Fetch извлекает, проверяет и возвращает конфигурацию
func (p configProvider) Fetch() (*config.Config, error) {
	c := config.NewConfig()
	for _, pr := range p.parsers {
		if err := pr.Parse(c); err != nil {
			return nil, fmt.Errorf("failed to parse multichecker config: %w", err)
		}
	}
	return c, nil
}

// New возвращает сервис configProvider для извлечения конфигурации multichecker
func New(parsers ...parser) *configProvider {
	return &configProvider{parsers}
}
