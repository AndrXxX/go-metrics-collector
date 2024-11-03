// # multichecker, состоящий из анализаторов:
//
// 1. Go анализаторы:
//   - printf.Analyzer
//   - shadow.Analyzer
//   - shift.Analyzer
//   - structtag.Analyzer
//
// 2. Дополнительные анализаторы:
//   - testifylint анализатор тестов https://github.com/Antonboom/testifylint
//   - errcheck анализатор необработанных ошибок https://github.com/kisielk/errcheck
//   - OSExitAnalyzer анализатор использования вызовов os.Exit() в функции main пакета main
//
// 3. Анализаторы staticcheck https://staticcheck.dev/
//
// По умолчанию включены анализаторы серии SA***, ST*** и QF***
//
// Можно указать другой набор анализаторов staticcheck:
//   - с помощью переменных env, например: STATIC_ANALYZERS=SA1000,SA1032,SA4004,QF.*,SA4005
//   - с помощью аргументов командной строки, например: -sa=SA1000,SA1032,SA4004,QF.*,SA4005
//
// Можно также исключить набор анализаторов staticcheck:
//   - с помощью переменных env, например: EXCLUDE_STATIC_ANALYZERS=SA1000,SA1032,SA4004,QF.*,SA4005
//   - с помощью аргументов командной строки, например: -esa=SA1000,SA1032,SA4004,QF.*,SA4005
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
