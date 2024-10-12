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

func (h *updateMetricsHandler) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Handle(w, r, nil)
	}
}

// Handle updates metrics from request
func (h *updateMetricsHandler) Handle(w http.ResponseWriter, r *http.Request, next http.Handler) {
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
	if next != nil {
		next.ServeHTTP(w, r)
	}
}

// New Return updateMetricsHandler
func New(u updater, f formatter, i identifier) *updateMetricsHandler {
	return &updateMetricsHandler{u, f, i}
}
