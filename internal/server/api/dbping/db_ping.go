package dbping

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
)

type dbPingHandler struct {
	c dbChecker
}

// Handler возвращает http.HandlerFunc
func (h *dbPingHandler) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h.c.Check(r.Context())
		if err != nil {
			logger.Log.Error("Error on check db", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// New возвращает экземпляр обработчика для проверки соединения с базой данных
func New(c dbChecker) *dbPingHandler {
	return &dbPingHandler{c}
}
