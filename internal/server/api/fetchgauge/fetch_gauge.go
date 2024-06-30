package fetchgauge

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type fetchGaugeHandler struct {
	s storage
}

func (h *fetchGaugeHandler) Handle(w http.ResponseWriter, r *http.Request) (ok bool) {
	metric := chi.URLParam(r, vars.Metric)
	if val, ok := h.s.Get(metric); ok {
		_, _ = fmt.Fprintf(w, "%v", *val.Value)
		w.WriteHeader(http.StatusOK)
		return true
	}
	w.WriteHeader(http.StatusNotFound)
	return false
}

func New(s storage) *fetchGaugeHandler {
	return &fetchGaugeHandler{s}
}
