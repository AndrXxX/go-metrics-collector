package dbstorage

import (
	"database/sql"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
	"go.uber.org/zap"
)

type dbStorage struct {
	db *sql.DB
}

func New(db *sql.DB) dbStorage {
	return dbStorage{db}
}

func (s *dbStorage) Insert(_ string, value *models.Metrics) {
	// TODO: realise with context
	stmt := `INSERT INTO metrics (name, type, delta, value) VALUES(?, ?, ?, ?)`
	_, err := s.db.Exec(stmt, value.ID, value.MType, value.Delta, value.Value)
	if err != nil {
		logger.Log.Error("Failed to insert metrics", zap.Error(err))
	}
}

func (s *dbStorage) Get(name string) (value *models.Metrics, ok bool) {
	// TODO: realise with context
	row := s.db.QueryRow("SELECT name, type, delta, value FROM metrics WHERE name = ?", name)
	v := models.Metrics{}
	err := row.Scan(&v.ID, &v.MType, &v.Delta, &v.Value)
	if err != nil {
		logger.Log.Error("error on scan all", zap.Error(err))
		return nil, false
	}
	return &v, false
}

func (s *dbStorage) All() map[string]*models.Metrics {
	list := make(map[string]*models.Metrics)

	// TODO: realise with context
	rows, err := s.db.Query("SELECT name, type, delta, value from metrics ")
	if err != nil {
		logger.Log.Error("error on select all", zap.Error(err))
		return list
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logger.Log.Error("close rows on all failed", zap.Error(err))
		}
	}(rows)

	for rows.Next() {
		v := models.Metrics{}
		err = rows.Scan(&v.ID, &v.MType, &v.Delta, &v.Value)
		if err != nil {
			logger.Log.Error("error on scan all", zap.Error(err))
			return list
		}
		list[v.ID] = &v
	}

	err = rows.Err()
	if err != nil {
		logger.Log.Error("error on fetch all", zap.Error(err))
	}
	return list
}
