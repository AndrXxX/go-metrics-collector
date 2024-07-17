package dbping

import (
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
	"go.uber.org/zap"
	"net/http"
)

type dbPingHandler struct {
	c dbChecker
}

func (h *dbPingHandler) Handle(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	err := h.c.Check(r.Context())
	if err != nil {
		logger.Log.Error("Error on check db", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if next != nil {
		next(w, r)
	}
}

func New(c dbChecker) *dbPingHandler {
	return &dbPingHandler{c}
}
