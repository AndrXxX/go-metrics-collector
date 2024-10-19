package configprovider

import (
	"fmt"

	"github.com/asaskevich/govalidator"

	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
)

type configProvider struct {
	parsers []parser
}

// Fetch извлекает, проверяет и возвращает конфигурацию
func (p configProvider) Fetch() (*config.Config, error) {
	settings := config.NewConfig()
	for _, pr := range p.parsers {
		if err := pr.Parse(settings); err != nil {
			return nil, fmt.Errorf("failed to parse config: %w", err)
		}
	}
	if _, err := govalidator.ValidateStruct(settings); err != nil {
		return nil, fmt.Errorf("failed to validate env vars: %w", err)
	}
	return settings, nil
}

// New возвращает сервис configProvider для извлечения конфигурации
func New(parsers ...parser) *configProvider {
	return &configProvider{parsers}
}
