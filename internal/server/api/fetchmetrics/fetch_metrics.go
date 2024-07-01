package fetchmetrics

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type fetchMetricsHandler struct {
	s           storage
	stringifier stringifier
}

func (h *fetchMetricsHandler) Handle(w http.ResponseWriter, r *http.Request) (ok bool) {
	metric := chi.URLParam(r, vars.Metric)
	val, ok := h.s.Get(metric)
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

func New(s storage, stringifier stringifier) *fetchMetricsHandler {
	return &fetchMetricsHandler{s, stringifier}
}
