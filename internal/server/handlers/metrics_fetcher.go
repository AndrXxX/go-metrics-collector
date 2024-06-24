package handlers

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/server/templates"
	"html/template"
	"log"
	"net/http"
)

type storage[T any] interface {
	All() map[string]T
}

func MetricsFetcher(gs storage[float64], cs storage[int64]) http.HandlerFunc {
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
			Items: fetchMetrics(gs, cs),
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

func fetchMetrics(gs storage[float64], cs storage[int64]) map[string]string {
	result := map[string]string{}
	for k, v := range cs.All() {
		result[k] = fmt.Sprintf("%d", v)
	}
	for k, v := range gs.All() {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result
}
