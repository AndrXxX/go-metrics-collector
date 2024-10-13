package updatemetrics

import (
	"fmt"
	"net/http"

	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
)

type updateMetricsHandler struct {
	u updater
	f formatter
	i identifier
}

// Handler returns HandlerFunc to update metrics from request
func (h *updateMetricsHandler) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metric, err := h.i.Process(r)
		if metric == nil || err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		metric, err = h.u.Update(r.Context(), metric)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		str, err := h.f.Format(metric)
		if err == nil {
			w.WriteHeader(http.StatusOK)
			_, err = fmt.Fprintf(w, "%s", str)
		}
		if err != nil {
			logger.Log.Error("Failed to write metrics response")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

// New Return updateMetricsHandler
func New(u updater, f formatter, i identifier) *updateMetricsHandler {
	return &updateMetricsHandler{u, f, i}
}
