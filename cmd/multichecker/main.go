package main

import (
	"fmt"
	"log"
	"regexp"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"

	"github.com/AndrXxX/go-metrics-collector/internal/multichecker/config"
	"github.com/AndrXxX/go-metrics-collector/internal/multichecker/services/configprovider"
)

func main() {
	c, err := configprovider.New().Fetch()
	if err != nil {
		log.Fatal(err)
	}

	var checks []*analysis.Analyzer
	checks = append(checks, getAnalysisChecks()...)

	staticChecks, err := getStaticChecks(c)
	if err != nil {
		log.Fatal(err)
	}
	checks = append(checks, staticChecks...)

	multichecker.Main(
		checks...,
	)
}

func getAnalysisChecks() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		printf.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		structtag.Analyzer,
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
