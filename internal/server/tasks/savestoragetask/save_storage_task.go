package savestoragetask

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
)

type saveStorageTask struct {
	i time.Duration
	s storageSaver
}

func (t *saveStorageTask) Execute(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := t.s.Save(ctx)
			if err != nil {
				logger.Log.Error("save storage task failed", zap.Error(err))
			}
		}
		time.Sleep(t.i)
	}
}

func New(i time.Duration, s storageSaver) *saveStorageTask {
	return &saveStorageTask{i, s}
}
