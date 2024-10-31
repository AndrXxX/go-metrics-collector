package main

import (
	"log"

	"golang.org/x/tools/go/analysis/multichecker"

	"github.com/AndrXxX/go-metrics-collector/internal/staticlint/services/analyzersprovider"
	"github.com/AndrXxX/go-metrics-collector/internal/staticlint/services/configprovider"
	"github.com/AndrXxX/go-metrics-collector/internal/staticlint/services/envparser"
)

func main() {
	c, err := configprovider.New(envparser.EnvParser{}).Fetch()
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
