package dbping

import (
	"database/sql"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
	"go.uber.org/zap"
	"net/http"
)

type dbPingHandler struct {
	db *sql.DB
}

func (h *dbPingHandler) Handle(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	err := h.db.Ping()
	if err != nil {
		logger.Log.Error("Error on ping db", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if next != nil {
		next(w, r)
	}
}

func New(db *sql.DB) *dbPingHandler {
	return &dbPingHandler{db}
}
