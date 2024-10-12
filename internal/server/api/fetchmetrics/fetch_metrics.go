package fetchmetrics

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
	"net/http"
)

type fetchMetricsHandler struct {
	s  storage[*models.Metrics]
	f  formatter
	i  identifier
	mc metricsChecker
}

func (h *fetchMetricsHandler) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Handle(w, r, nil)
	}
}

func (h *fetchMetricsHandler) Handle(w http.ResponseWriter, r *http.Request, next http.Handler) {
	metric, err := h.i.Process(r)
	if metric == nil || err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if !h.mc.IsValid(metric) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	val, ok := h.s.Get(r.Context(), metric.ID)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	str, err := h.f.Format(val)
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

func New(s storage[*models.Metrics], f formatter, i identifier, mc metricsChecker) *fetchMetricsHandler {
	return &fetchMetricsHandler{s, f, i, mc}
}
