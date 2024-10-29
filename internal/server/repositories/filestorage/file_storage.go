package filestorage

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"github.com/AndrXxX/go-metrics-collector/internal/server/repositories"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/storagesaver"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
)

type fileStorage struct {
	c  *config.Config
	s  repositories.Storage[*models.Metrics]
	ss storageSaver
}

// New возвращает хранилище метрик в файле
func New(c *config.Config, s repositories.Storage[*models.Metrics]) fileStorage {
	ss := storagesaver.New(c.FileStoragePath, s, c.RepeatIntervals)
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

// Insert вставляет запись
func (s *fileStorage) Insert(ctx context.Context, name string, value *models.Metrics) {
	s.s.Insert(ctx, name, value)
}

// Get извлекает запись
func (s *fileStorage) Get(ctx context.Context, name string) (value *models.Metrics, ok bool) {
	val, found := s.s.Get(ctx, name)
	return val, found
}

// All извлекает все записи
func (s *fileStorage) All(ctx context.Context) map[string]*models.Metrics {
	return s.s.All(ctx)
}

// Delete удаляет запись
func (s *fileStorage) Delete(ctx context.Context, name string) (ok bool) {
	return s.s.Delete(ctx, name)
}

// Shutdown завершение работы хранилища
func (s *fileStorage) Shutdown(ctx context.Context) error {
	return s.Save(ctx)
}

// Save сохранение хранилища
func (s *fileStorage) Save(ctx context.Context) error {
	err := s.ss.Save(ctx)
	if err != nil {
		return fmt.Errorf("error saving storage: %w", err)
	}
	return nil
}
