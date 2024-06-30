package fetchmetrics

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"github.com/AndrXxX/go-metrics-collector/internal/server/templates"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
	"html/template"
	"net/http"
)

type fetchMetricsHandler struct {
	gs storage[*models.Metrics]
	cs storage[*models.Metrics]
}

func (h *fetchMetricsHandler) Handle(w http.ResponseWriter, _ *http.Request) (ok bool) {
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
		result[k] = fmt.Sprintf("%d", *v.Delta)
	}
	for k, v := range h.gs.All() {
		result[k] = fmt.Sprintf("%v", *v.Value)
	}
	return result
}

func New(gs storage[*models.Metrics], cs storage[*models.Metrics]) *fetchMetricsHandler {
	return &fetchMetricsHandler{gs, cs}
}
