package fetchcounter

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/logger"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type fetchCounterHandler struct {
	s cfStorage
}

func (h *fetchCounterHandler) Handle(w http.ResponseWriter, r *http.Request) (ok bool) {
	w.Header().Set("Content-Type", "text/plain")
	metric := chi.URLParam(r, vars.Metric)
	val, ok := h.s.Get(metric)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return false
	}
	_, err := w.Write([]byte(fmt.Sprintf("%d", val)))
	if err != nil {
		logger.Log.Error("Failed to write counter response")
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}
	w.WriteHeader(http.StatusOK)
	return true
}

func New(s cfStorage) *fetchCounterHandler {
	return &fetchCounterHandler{s}
}
