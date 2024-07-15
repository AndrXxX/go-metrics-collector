package dbstorage

import (
	"context"
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

func (s *dbStorage) Insert(ctx context.Context, name string, value *models.Metrics) {
	stmt := `INSERT INTO metrics (name, type, delta, value) VALUES($1, $2, $3, $4)`
	_, err := s.db.ExecContext(ctx, stmt, name, value.MType, value.Delta, value.Value)
	if err != nil {
		logger.Log.Error("Failed to insert metrics", zap.Error(err))
	}
}

func (s *dbStorage) Get(ctx context.Context, name string) (value *models.Metrics, ok bool) {
	row := s.db.QueryRowContext(ctx, "SELECT name, type, delta, value FROM metrics WHERE name = $1", name)
	v := models.Metrics{}
	err := row.Scan(&v.ID, &v.MType, &v.Delta, &v.Value)
	if err != nil {
		return nil, false
	}
	return &v, true
}

func (s *dbStorage) All(ctx context.Context) map[string]*models.Metrics {
	list := make(map[string]*models.Metrics)

	rows, err := s.db.QueryContext(ctx, "SELECT name, type, delta, value from metrics ")
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
