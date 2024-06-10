package handlers

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/repositories"
	"html/template"
	"log"
	"net/http"
)

const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
	</head>
	<body>
		<h1>{{.Title}}</h1>
		{{ range $metric, $value := .Items }}
			<div><strong>{{ $metric }}:</strong> <span>{{ $value }}</span></div>
		{{else}}
			<div><strong>Список пуст</strong></div>
		{{end}}
	</body>
</html>`

func MetricsFetcher(s repositories.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		t, err := template.New("webpage").Parse(tpl)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		data := struct {
			Title string
			Items map[string]string
		}{
			Title: "Metrics List",
			Items: fetchMetrics(s),
		}

		err = t.Execute(w, data)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func fetchMetrics(s repositories.Repository) map[string]string {
	result := map[string]string{}
	for k, v := range s.GetCounterAll() {
		result[k] = fmt.Sprintf("%d", v)
	}
	for k, v := range s.GetGaugeAll() {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result
}
