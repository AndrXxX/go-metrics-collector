package filestorage

import (
	"context"
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"github.com/AndrXxX/go-metrics-collector/internal/server/repositories"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/storagesaver"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
	"go.uber.org/zap"
)

type fileStorage struct {
	c  *config.Config
	s  repositories.Storage[*models.Metrics]
	ss storageSaver
}

func New(c *config.Config, s repositories.Storage[*models.Metrics]) fileStorage {
	ss := storagesaver.New(c.FileStoragePath, s)
	if c.Restore {
		err := ss.Restore(context.TODO())
		if err != nil {
			logger.Log.Error("Error restoring storage", zap.Error(err))
		}
	}
	return fileStorage{
		c,
		s,
		ss,
	}
}

func (s *fileStorage) Insert(ctx context.Context, name string, value *models.Metrics) {
	s.s.Insert(ctx, name, value)
}

func (s *fileStorage) Get(ctx context.Context, name string) (value *models.Metrics, ok bool) {
	val, found := s.s.Get(ctx, name)
	return val, found
}

func (s *fileStorage) All(ctx context.Context) map[string]*models.Metrics {
	return s.s.All(ctx)
}

func (s *fileStorage) Shutdown(ctx context.Context) error {
	return s.Save(ctx)
}

func (s *fileStorage) Save(ctx context.Context) error {
	err := s.ss.Save(ctx)
	if err != nil {
		return fmt.Errorf("error saving storage: %w", err)
	}
	return nil
}
