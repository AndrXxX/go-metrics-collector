package honnef

import (
	"fmt"

	"golang.org/x/tools/go/analysis"
	"honnef.co/go/tools/analysis/lint"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"

	"github.com/AndrXxX/go-metrics-collector/internal/staticlint/config"
	"github.com/AndrXxX/go-metrics-collector/internal/staticlint/services/checksprovider/filters"
)

// Analyzers возвращает список анализаторов honnef на основе конфигурации
func Analyzers(c *config.Config) ([]*analysis.Analyzer, error) {
	var checks []*analysis.Analyzer
	raw := [][]*analysis.Analyzer{
		convert(staticcheck.Analyzers),
		convert(quickfix.Analyzers),
		convert(stylecheck.Analyzers),
		convert(simple.Analyzers),
	}
	for _, pack := range raw {
		filtered, err := filters.ByName(pack, c.StaticChecks)
		if err != nil {
			return nil, fmt.Errorf("failed to filter checks: %v", err)
		}
		checks = append(checks, filtered...)
	}
	return checks, nil
}

func convert(list []*lint.Analyzer) []*analysis.Analyzer {
	converted := make([]*analysis.Analyzer, len(list))
	for i, v := range list {
		converted[i] = v.Analyzer
	}
	return converted
}
