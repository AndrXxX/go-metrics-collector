package flagsparser

import (
	fl "flag"
	"strings"

	"github.com/AndrXxX/go-metrics-collector/internal/staticlint/config"
)

// FlagsParser сервис для парсинга аргументов командной строки для staticlint
type FlagsParser struct {
}

// Parse парсит аргументы командной строки и наполняет конфигурацию для staticlint
func (p FlagsParser) Parse(c *config.Config) error {
	var staticAnalyzers string
	var excludeStaticAnalyzers string
	fl.StringVar(&staticAnalyzers, "sa", strings.Join(c.StaticAnalyzers, ","), "pass static analyzers names separated by comma")
	fl.StringVar(&excludeStaticAnalyzers, "esa", strings.Join(c.ExcludeStaticAnalyzers, ","), "pass static analyzers names separated by comma to exclude")
	fl.Parse()
	if staticAnalyzers != "" {
		c.StaticAnalyzers = strings.Split(staticAnalyzers, ",")
	}
	if excludeStaticAnalyzers != "" {
		c.ExcludeStaticAnalyzers = strings.Split(excludeStaticAnalyzers, ",")
	}
	return nil
}
