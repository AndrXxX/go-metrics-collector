package checksprovider

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
	"github.com/AndrXxX/go-metrics-collector/internal/staticlint/services/checksprovider/honnef"
	"github.com/AndrXxX/go-metrics-collector/internal/staticlint/services/osexitanalyzer"
)

// ChecksProvider возвращает сервис для сбора списка проверок
type ChecksProvider struct {
}

// Fetch возвращает список проверок на основе конфигурации
func (p ChecksProvider) Fetch(c *config.Config) ([]*analysis.Analyzer, error) {
	var checks []*analysis.Analyzer
	checks = append(checks, getAnalysisChecks()...)
	checks = append(checks, getAdditionalChecks()...)

	staticChecks, err := honnef.Analyzers(c)
	if err != nil {
		return nil, fmt.Errorf("error on getting static checks: %w", err)
	}
	checks = append(checks, staticChecks...)
	return checks, nil
}

func getAnalysisChecks() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		printf.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		structtag.Analyzer,
	}
}

func getAdditionalChecks() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		testifyAnalyzer.New(),
		errcheck.Analyzer,
		osexitanalyzer.OSExitAnalyzer,
	}
}
