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
	if !exist {
		currentModel = newModel
		u.s.Insert(ctx, currentModel.ID, currentModel)
	}
	if newModel.MType == metrics.Gauge {
		currentModel.Value = newModel.Value
		return currentModel, nil
	}
	if newModel.Delta != nil && currentModel.Delta != nil && exist {
		newVal := *currentModel.Delta + *newModel.Delta
		currentModel.Delta = &newVal
		return currentModel, nil
	}
	currentModel.Delta = newModel.Delta
	return currentModel, nil
}
