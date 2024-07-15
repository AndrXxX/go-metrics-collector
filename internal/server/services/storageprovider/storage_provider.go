package storageprovider

import (
	"database/sql"
	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
	"github.com/AndrXxX/go-metrics-collector/internal/server/interfaces"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"github.com/AndrXxX/go-metrics-collector/internal/server/repositories/dbstorage"
	"github.com/AndrXxX/go-metrics-collector/internal/server/repositories/filestorage"
	"github.com/AndrXxX/go-metrics-collector/internal/server/repositories/memory"
)

type storageProvider struct {
	c  *config.Config
	db *sql.DB
}

func New(c *config.Config, db *sql.DB) *storageProvider {
	return &storageProvider{c, db}
}

func (sp *storageProvider) Storage() interfaces.MetricsStorage {
	if sp.c.DatabaseDSN != "" {
		s := dbstorage.New(sp.db)
		return &s
	}
	if sp.c.FileStoragePath != "" {
		ms := memory.New[*models.Metrics]()
		s := filestorage.New(sp.c, &ms)
		return &s
	}
	s := memory.New[*models.Metrics]()
	return &s
}
