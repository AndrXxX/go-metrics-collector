package dbstorage

import (
	"database/sql"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
)

type dbStorage struct {
	db *sql.DB
}

func New(db *sql.DB) dbStorage {
	return dbStorage{db}
}

func (s *dbStorage) Insert(name string, value *models.Metrics) {
	// TODO
}

func (s *dbStorage) Get(name string) (value *models.Metrics, ok bool) {
	// TODO
	return &models.Metrics{}, false
}

func (s *dbStorage) All() map[string]*models.Metrics {
	// TODO
	return make(map[string]*models.Metrics)
}
