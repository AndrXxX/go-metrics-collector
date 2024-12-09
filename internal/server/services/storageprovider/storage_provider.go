package storageprovider

import (
	"context"
	"database/sql"
	"time"

	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
	"github.com/AndrXxX/go-metrics-collector/internal/server/interfaces"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"github.com/AndrXxX/go-metrics-collector/internal/server/repositories/dbstorage"
	"github.com/AndrXxX/go-metrics-collector/internal/server/repositories/filestorage"
	"github.com/AndrXxX/go-metrics-collector/internal/server/repositories/memory"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/storagesaver"
	"github.com/AndrXxX/go-metrics-collector/internal/server/tasks/savestoragetask"
)

type storageProvider struct {
	c  *config.Config
	db *sql.DB
}

// New возвращает сервис для предоставления хранилища метрик
func New(c *config.Config, db *sql.DB) *storageProvider {
	return &storageProvider{c, db}
}

// Storage возвращает хранилище метрик
func (sp *storageProvider) Storage(ctx context.Context) interfaces.MetricsStorage {
	if sp.c.DatabaseDSN != "" {
		s := dbstorage.New(sp.db, sp.c.RepeatIntervals)
		return &s
	}
	ms := memory.New[*models.Metrics]()
	if sp.c.FileStoragePath != "" {
		s := filestorage.New(&ms, filestorage.WithStorageSaver(sp.c, storagesaver.New(sp.c.FileStoragePath, &ms, sp.c.RepeatIntervals)))
		sst := savestoragetask.New(time.Duration(sp.c.StoreInterval)*time.Second, s)
		go sst.Execute(ctx)
		return s
	}
	return &ms
}
