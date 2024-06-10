package handlers

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/repositories"
	"github.com/AndrXxX/go-metrics-collector/internal/templates"
	"html/template"
	"log"
	"net/http"
)

func MetricsFetcher(s repositories.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		t, err := template.New("webpage").Parse(templates.MetricsList)
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
