package fetchcounter

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type fetchCounterHandler struct {
	s storage
}

func (h *fetchCounterHandler) Handle(w http.ResponseWriter, r *http.Request) (ok bool) {
	metric := chi.URLParam(r, vars.Metric)
	val, ok := h.s.Get(metric)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return false
	}
	_, err := fmt.Fprintf(w, "%d", *val.Delta)
	if err != nil {
		logger.Log.Error("Failed to write counter response")
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}
	w.WriteHeader(http.StatusOK)
	return true
}

func New(s storage) *fetchCounterHandler {
	return &fetchCounterHandler{s}
}
