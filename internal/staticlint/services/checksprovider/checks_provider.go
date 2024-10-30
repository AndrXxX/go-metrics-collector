package checksprovider

import (
	"fmt"
	"regexp"

	testifyAnalyzer "github.com/Antonboom/testifylint/analyzer"
	"github.com/kisielk/errcheck/errcheck"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"

	"github.com/AndrXxX/go-metrics-collector/internal/staticlint/config"
)

// ChecksProvider возвращает сервис для сбора списка проверок
type ChecksProvider struct {
}

// Fetch возвращает список проверок на основе конфигурации
func (p ChecksProvider) Fetch(c *config.Config) ([]*analysis.Analyzer, error) {
	var checks []*analysis.Analyzer
	checks = append(checks, getAnalysisChecks()...)
	checks = append(checks, getAdditionalChecks()...)

	staticChecks, err := getStaticChecks(c)
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
	}
}

func getStaticChecks(c *config.Config) ([]*analysis.Analyzer, error) {
	var checks []*analysis.Analyzer
	for _, v := range staticcheck.Analyzers {
		for _, check := range c.StaticChecks {
			matched, err := regexp.MatchString(check, v.Analyzer.Name)
			if err != nil {
				return nil, fmt.Errorf("failed to match static check '%s': %v", v.Analyzer.Name, err)
			}
			if matched {
				checks = append(checks, v.Analyzer)
			}
		}
	}
	return checks, nil
}
