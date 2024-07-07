package updatemetrics

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
	"net/http"
)

type updateMetricsHandler struct {
	u updater
	f formatter
	i identifier
}

// Handle updates metrics from request
func (h *updateMetricsHandler) Handle(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	metric, err := h.i.Process(r)
	if metric == nil || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	metric, err = h.u.Update(metric)
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
	if next != nil {
		next(w, r)
	}
}

// New Return updateMetricsHandler
func New(u updater, f formatter, i identifier) *updateMetricsHandler {
	return &updateMetricsHandler{u, f, i}
}
