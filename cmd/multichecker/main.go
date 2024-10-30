package main

import (
	"log"

	"golang.org/x/tools/go/analysis/multichecker"

	"github.com/AndrXxX/go-metrics-collector/internal/multichecker/services/checksprovider"
	"github.com/AndrXxX/go-metrics-collector/internal/multichecker/services/configprovider"
)

func main() {
	c, err := configprovider.New().Fetch()
	if err != nil {
		log.Fatal(err)
	}
	checks, err := checksprovider.ChecksProvider{}.Fetch(c)
	if err != nil {
		log.Fatal(err)
	}
	multichecker.Main(
		checks...,
	)
}
