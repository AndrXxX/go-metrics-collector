package filestorage

import (
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
		err := ss.Restore()
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

func (s *fileStorage) Insert(name string, value *models.Metrics) {
	s.s.Insert(name, value)
}

func (s *fileStorage) Get(name string) (value *models.Metrics, ok bool) {
	val, found := s.s.Get(name)
	return val, found
}

func (s *fileStorage) All() map[string]*models.Metrics {
	return s.s.All()
}

func (s *fileStorage) Shutdown() error {
	return s.Save()
}

func (s *fileStorage) Save() error {
	err := s.ss.Save()
	if err != nil {
		return fmt.Errorf("error saving storage: %w", err)
	}
	return nil
}
