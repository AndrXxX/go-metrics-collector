package configprovider

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
	"github.com/asaskevich/govalidator"
)

type configProvider struct {
	fp flagsParser
	ep envParser
}

func (p configProvider) Fetch() (*config.Config, error) {
	settings := config.NewConfig()
	p.fp.Parse(settings)
	if err := p.ep.Parse(settings); err != nil {
		return nil, fmt.Errorf("failed to parse env vars: %w", err)
	}
	if _, err := govalidator.ValidateStruct(settings); err != nil {
		return nil, fmt.Errorf("failed to validate env vars: %w", err)
	}
	return settings, nil
}

func New(fp flagsParser, ep envParser) *configProvider {
	return &configProvider{fp, ep}
}
