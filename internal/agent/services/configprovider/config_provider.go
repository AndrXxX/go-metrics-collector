package configprovider

import (
	"fmt"

	"github.com/asaskevich/govalidator"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
)

type configProvider struct {
	parsers []parser
}

// Fetch извлекает, проверяет и возвращает конфигурацию
func (p configProvider) Fetch() (*config.Config, error) {
	c := config.NewConfig()
	for _, pr := range p.parsers {
		if err := pr.Parse(c); err != nil {
			return nil, fmt.Errorf("failed to parse agent config: %w", err)
		}
	}
	if _, err := govalidator.ValidateStruct(c.Intervals); err != nil {
		return nil, fmt.Errorf("failed to validate Intervals config: %w", err)
	}
	if _, err := govalidator.ValidateStruct(c.Common); err != nil {
		return nil, fmt.Errorf("failed to validate Common config: %w", err)
	}
	return c, nil
}

// New возвращает сервис configProvider для извлечения конфигурации
func New(parsers ...parser) *configProvider {
	return &configProvider{parsers}
}
