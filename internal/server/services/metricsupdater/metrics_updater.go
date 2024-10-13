package metricsupdater

import (
	"context"

	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
)

type metricsUpdater struct {
	s storage[*models.Metrics]
}

func New(s storage[*models.Metrics]) *metricsUpdater {
	return &metricsUpdater{s}
}

func (u *metricsUpdater) Update(ctx context.Context, newModel *models.Metrics) (*models.Metrics, error) {
	currentModel, exist := u.s.Get(ctx, newModel.ID)
	if exist {
		u.s.Delete(ctx, newModel.ID)
	}
	if exist && newModel.MType == metrics.Counter {
		newVal := *currentModel.Delta + *newModel.Delta
		newModel.Delta = &newVal
	}
	u.s.Insert(ctx, newModel.ID, newModel)
	return currentModel, nil
}

func (u *metricsUpdater) UpdateMany(ctx context.Context, list []models.Metrics) error {
	for _, model := range list {
		_, err := u.Update(ctx, &model)
		if err != nil {
			return err
		}
	}
	return nil
}
