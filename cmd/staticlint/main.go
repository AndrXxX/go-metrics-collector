package main

import (
	"log"

	"golang.org/x/tools/go/analysis/multichecker"

	"github.com/AndrXxX/go-metrics-collector/internal/staticlint/services/analyzersprovider"
	"github.com/AndrXxX/go-metrics-collector/internal/staticlint/services/configprovider"
	"github.com/AndrXxX/go-metrics-collector/internal/staticlint/services/envparser"
	"github.com/AndrXxX/go-metrics-collector/internal/staticlint/services/flagsparser"
)

func main() {
	c, err := configprovider.New(envparser.EnvParser{}, flagsparser.FlagsParser{}).Fetch()
	if err != nil {
		log.Fatal(err)
	}
	analyzers, err := analyzersprovider.AnalyzersProvider{}.Fetch(c)
	if err != nil {
		log.Fatal(err)
	}
	multichecker.Main(
		analyzers...,
	)
}
