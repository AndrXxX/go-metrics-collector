package analyzersprovider

import (
	"fmt"

	testifyAnalyzer "github.com/Antonboom/testifylint/analyzer"
	"github.com/kisielk/errcheck/errcheck"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"

	"github.com/AndrXxX/go-metrics-collector/internal/staticlint/config"
	"github.com/AndrXxX/go-metrics-collector/internal/staticlint/services/analyzersprovider/honnef"
	"github.com/AndrXxX/go-metrics-collector/internal/staticlint/services/osexitanalyzer"
)

// AnalyzersProvider возвращает сервис для сбора списка проверок
type AnalyzersProvider struct {
}

// Fetch возвращает список анализаторов на основе конфигурации
func (p AnalyzersProvider) Fetch(c *config.Config) ([]*analysis.Analyzer, error) {
	var list []*analysis.Analyzer
	list = append(list, getGolangAnalyzers()...)
	list = append(list, getAdditionalAnalyzers()...)

	staticAnalyzers, err := honnef.Analyzers(c)
	if err != nil {
		return nil, fmt.Errorf("error on getting static analyzers: %w", err)
	}
	list = append(list, staticAnalyzers...)
	return list, nil
}

func getGolangAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		printf.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		structtag.Analyzer,
	}
}

func getAdditionalAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		testifyAnalyzer.New(),
		errcheck.Analyzer,
		osexitanalyzer.OSExitAnalyzer,
	}
}
