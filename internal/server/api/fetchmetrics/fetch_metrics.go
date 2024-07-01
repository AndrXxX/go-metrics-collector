package fetchmetrics

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
	"net/http"
)

type fetchMetricsHandler struct {
	s           storage
	stringifier stringifier
	i           identifier
}

func (h *fetchMetricsHandler) Handle(w http.ResponseWriter, r *http.Request) (ok bool) {
	metric, err := h.i.Process(r)
	if metric == nil || err != nil {
		w.WriteHeader(http.StatusNotFound)
		return false
	}
	val, ok := h.s.Get(metric.ID)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return false
	}
	str, err := h.stringifier.String(val)
	if err == nil {
		_, err = fmt.Fprintf(w, "%s", str)
	}
	if err != nil {
		logger.Log.Error("Failed to write metrics response")
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}
	w.WriteHeader(http.StatusOK)
	return true
}

func New(s storage, stringifier stringifier, i identifier) *fetchMetricsHandler {
	return &fetchMetricsHandler{s, stringifier, i}
}
