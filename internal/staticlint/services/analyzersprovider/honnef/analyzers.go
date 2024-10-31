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
	"github.com/AndrXxX/go-metrics-collector/internal/staticlint/services/analyzersprovider/filters"
)

// Analyzers возвращает список анализаторов honnef на основе конфигурации
func Analyzers(c *config.Config) ([]*analysis.Analyzer, error) {
	var list []*analysis.Analyzer
	raw := [][]*analysis.Analyzer{
		convert(staticcheck.Analyzers),
		convert(quickfix.Analyzers),
		convert(stylecheck.Analyzers),
		convert(simple.Analyzers),
	}
	for _, pack := range raw {
		filtered, err := filters.ByName(pack, c.StaticAnalyzers)
		if err != nil {
			return nil, fmt.Errorf("failed to filter analyzers: %v", err)
		}
		list = append(list, filtered...)
	}
	return list, nil
}

func convert(list []*lint.Analyzer) []*analysis.Analyzer {
	converted := make([]*analysis.Analyzer, len(list))
	for i, v := range list {
		converted[i] = v.Analyzer
	}
	return converted
}
