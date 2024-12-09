package filestorage

import (
	"context"

	"go.uber.org/zap"

	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
)

func WithStorageSaver(c *config.Config, ss storageSaver) Option {
	return func(fs *fileStorage) {
		fs.ss = ss
		if c.Restore {
			err := ss.Restore(context.Background())
			if err != nil {
				logger.Log.Error("Error restoring storage", zap.Error(err))
			}
		}
	}
}
