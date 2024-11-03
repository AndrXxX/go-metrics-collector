package dbstorage

import (
	"context"
	"database/sql"
	"time"

	"go.uber.org/zap"

	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
)

type dbStorage struct {
	db *sql.DB
	ri []int
}

// New возвращает хранилище метрик в базе данных
func New(db *sql.DB, repeatIntervals []int) dbStorage {
	return dbStorage{db, repeatIntervals}
}

// Insert вставляет запись
func (s *dbStorage) Insert(ctx context.Context, name string, value *models.Metrics) {
	stmt, err := s.db.PrepareContext(ctx, "INSERT INTO metrics (name, type, delta, value) VALUES($1, $2, $3, $4)")
	if err != nil {
		logger.Log.Error("Failed to insert metrics", zap.Error(err))
		return
	}
	defer func(stmt *sql.Stmt) {
		if cErr := stmt.Close(); cErr != nil {
			logger.Log.Error("error on close stmt", zap.Error(cErr))
		}
	}(stmt)

	op := func() error {
		_, eRrr := stmt.ExecContext(ctx, name, value.MType, value.Delta, value.Value)
		return eRrr
	}
	err = op()
	if err != nil {
		err = s.repeat(op)
	}
	if err != nil {
		logger.Log.Error("Failed to insert metrics", zap.Error(err))
	}
}

// Get извлекает запись
func (s *dbStorage) Get(ctx context.Context, name string) (value *models.Metrics, ok bool) {
	row := s.db.QueryRowContext(ctx, "SELECT name, type, delta, value FROM metrics WHERE name = $1", name)
	v := models.Metrics{}
	err := row.Scan(&v.ID, &v.MType, &v.Delta, &v.Value)
	if err != nil {
		return nil, false
	}
	return &v, true
}

// All извлекает все записи
func (s *dbStorage) All(ctx context.Context) map[string]*models.Metrics {
	list := make(map[string]*models.Metrics)

	rows, err := s.db.QueryContext(ctx, "SELECT name, type, delta, value from metrics ")
	if err != nil {
		logger.Log.Error("error on select all", zap.Error(err))
		return list
	}
	defer func(rows *sql.Rows) {
		if cErr := rows.Close(); cErr != nil {
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

// Delete удаляет запись
func (s *dbStorage) Delete(ctx context.Context, name string) (ok bool) {
	stmt := `DELETE FROM metrics WHERE name = $1`
	_, err := s.db.ExecContext(ctx, stmt, name)
	if err != nil {
		logger.Log.Error("Failed to delete metrics", zap.Error(err))
		return false
	}
	return true
}

func (s *dbStorage) repeat(f func() error) (err error) {
	for _, repeatInterval := range s.ri {
		time.Sleep(time.Duration(repeatInterval) * time.Second)
		err = f()
		if err == nil {
			return nil
		}
	}
	return err
}
