package fetchmetrics

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/logger"
	"github.com/AndrXxX/go-metrics-collector/internal/server/templates"
	"html/template"
	"net/http"
)

type fetchMetricsHandler struct {
	gs mfStorage[float64]
	cs mfStorage[int64]
}

func (h *fetchMetricsHandler) Handle(w http.ResponseWriter, r *http.Request) (ok bool) {
	w.Header().Set("Content-Type", "text/html")
	t, err := template.New("webpage").Parse(templates.MetricsList)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Error on parse MetricsList template: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}

	data := struct {
		Title string
		Items map[string]string
	}{
		Title: "Metrics List",
		Items: h.fetchMetrics(),
	}

	err = t.Execute(w, data)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Error on execute MetricsList template: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}
	w.WriteHeader(http.StatusOK)
	return true
}

func (h *fetchMetricsHandler) fetchMetrics() map[string]string {
	result := map[string]string{}
	for k, v := range h.cs.All() {
		result[k] = fmt.Sprintf("%d", v)
	}
	for k, v := range h.gs.All() {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result
}

func New(gs mfStorage[float64], cs mfStorage[int64]) *fetchMetricsHandler {
	return &fetchMetricsHandler{gs, cs}
}
