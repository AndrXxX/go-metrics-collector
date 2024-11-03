package config

import "github.com/AndrXxX/go-metrics-collector/internal/staticlint/vars"

// NewConfig предоставляет конфигурацию со значениями по умолчанию для multichecker
func NewConfig() *Config {
	return &Config{
		StaticAnalyzers: []string{
			vars.StaticSAAnalyzers,
			vars.StaticSTAnalyzers,
			vars.StaticQFAnalyzers,
		},
		ExcludeStaticAnalyzers: []string{
			vars.StaticST1000Analyzer,
		},
	}
}
