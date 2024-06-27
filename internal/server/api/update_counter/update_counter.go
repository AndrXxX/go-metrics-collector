package update_counter

import (
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type updateCounterHandler struct {
	u updater
}

func (cu *updateCounterHandler) Handle(w http.ResponseWriter, r *http.Request) (ok bool) {
	w.Header().Set("Content-Type", "text/plain")
	metric := chi.URLParam(r, vars.Metric)
	value := chi.URLParam(r, vars.Value)
	if converted, err := strconv.ParseInt(value, 10, 64); err == nil {
		cu.u.Update(metric, converted)
		w.WriteHeader(http.StatusOK)
		return true
	}
	w.WriteHeader(http.StatusBadRequest)
	return false
}

func New(u updater) *updateCounterHandler {
	return &updateCounterHandler{u}
}
