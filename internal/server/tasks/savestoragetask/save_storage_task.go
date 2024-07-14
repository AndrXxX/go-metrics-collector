package savestoragetask

import (
	"context"
	"github.com/AndrXxX/go-metrics-collector/internal/server/repositories"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
	"go.uber.org/zap"
	"time"
)

type saveStorageTask struct {
	i time.Duration
	s repositories.StorageSaver
}

func (t *saveStorageTask) Execute(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := t.s.Save()
			if err != nil {
				logger.Log.Error("save storage task failed", zap.Error(err))
			}
		}
		time.Sleep(t.i)
	}
}

func New(i time.Duration, s repositories.StorageSaver) *saveStorageTask {
	return &saveStorageTask{i, s}
}
