package fetchmetrics

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
	"net/http"
)

type fetchMetricsHandler struct {
	sp storageProvider
	s  stringifier
	i  identifier
}

func (h *fetchMetricsHandler) Handle(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	metric, err := h.i.Process(r)
	if metric == nil || err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	storage := h.sp.GetStorage(metric.MType)
	if storage == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	val, ok := storage.Get(metric.ID)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	str, err := h.s.String(val)
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

func New(sp storageProvider, s stringifier, i identifier) *fetchMetricsHandler {
	return &fetchMetricsHandler{sp, s, i}
}
