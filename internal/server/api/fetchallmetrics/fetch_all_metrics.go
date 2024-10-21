package fetchallmetrics

import (
	"fmt"
	"html/template"
	"net/http"

	"go.uber.org/zap"

	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"github.com/AndrXxX/go-metrics-collector/internal/server/templates"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
)

type fetchAllMetricsHandler struct {
	s storage[*models.Metrics]
}

// Handler возвращает http.HandlerFunc
func (h *fetchAllMetricsHandler) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.New("webpage").Parse(templates.MetricsList)
		if err != nil {
			logger.Log.Error("Error on parse MetricsList template", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		data := struct {
			Title string
			Items map[string]string
		}{
			Title: "Metrics List",
			Items: h.fetchMetrics(r),
		}

		w.WriteHeader(http.StatusOK)
		err = t.Execute(w, data)
		if err != nil {
			logger.Log.Error("Error on execute MetricsList template", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (h *fetchAllMetricsHandler) fetchMetrics(r *http.Request) map[string]string {
	result := map[string]string{}
	for k, v := range h.s.All(r.Context()) {
		if v.MType == metrics.Counter {
			result[k] = fmt.Sprintf("%d", *v.Delta)
		} else {
			result[k] = fmt.Sprintf("%v", *v.Value)
		}
	}
	return result
}

// New возвращает обработчик для получения страницы со всеми метриками
func New(s storage[*models.Metrics]) *fetchAllMetricsHandler {
	return &fetchAllMetricsHandler{s}
}
