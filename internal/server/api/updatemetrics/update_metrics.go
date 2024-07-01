package updatemetrics

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
	"net/http"
)

type updateMetricsHandler struct {
	u updater
	s stringifier
	i identifier
}

func (h *updateMetricsHandler) Handle(w http.ResponseWriter, r *http.Request) (ok bool) {
	metric, err := h.i.Process(r)
	if metric == nil || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	metric, err = h.u.Update(metric)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	str, err := h.s.String(metric)
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

func New(u updater, s stringifier, i identifier) *updateMetricsHandler {
	return &updateMetricsHandler{u, s, i}
}
