package fetchallmetrics

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"github.com/AndrXxX/go-metrics-collector/internal/server/templates"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
	"html/template"
	"net/http"
)

type fetchMetricsHandler struct {
	s storage[*models.Metrics]
}

func (h *fetchMetricsHandler) Handle(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
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
		Items: h.fetchMetrics(),
	}

	w.WriteHeader(http.StatusOK)
	err = t.Execute(w, data)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Error on execute MetricsList template: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if next != nil {
		next(w, r)
	}
}

func (h *fetchMetricsHandler) fetchMetrics() map[string]string {
	result := map[string]string{}
	for k, v := range h.s.All() {
		if v.MType == metrics.Counter {
			result[k] = fmt.Sprintf("%d", *v.Delta)
		} else {
			result[k] = fmt.Sprintf("%v", *v.Value)
		}
	}
	return result
}

func New(s storage[*models.Metrics]) *fetchMetricsHandler {
	return &fetchMetricsHandler{s}
}
