package fetchmetrics

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/logger"
	"github.com/AndrXxX/go-metrics-collector/internal/server/templates"
	"html/template"
	"net/http"
)

type mfStorage[T any] interface {
	All() map[string]T
}

func MetricsFetcher(gs mfStorage[float64], cs mfStorage[int64]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		t, err := template.New("webpage").Parse(templates.MetricsList)
		if err != nil {
			logger.Log.Error(fmt.Sprintf("Error on parse MetricsList template: %s", err.Error()))
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
			logger.Log.Error(fmt.Sprintf("Error on execute MetricsList template: %s", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func fetchMetrics(gs mfStorage[float64], cs mfStorage[int64]) map[string]string {
	result := map[string]string{}
	for k, v := range cs.All() {
		result[k] = fmt.Sprintf("%d", v)
	}
	for k, v := range gs.All() {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result
}
