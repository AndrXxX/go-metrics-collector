package metricsupdater

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
)

type metricsUpdater struct {
	sp storageProvider
}

func New(sp storageProvider) *metricsUpdater {
	return &metricsUpdater{sp}
}

func (u *metricsUpdater) Update(newModel *models.Metrics) (*models.Metrics, error) {
	storage := u.sp.GetStorage(newModel.MType)
	if storage == nil {
		return nil, fmt.Errorf("not found storage for metrics type %s", newModel.MType)
	}
	currentModel, exist := storage.Get(newModel.ID)
	if !exist {
		currentModel = newModel
		storage.Insert(currentModel.ID, currentModel)
	}
	if newModel.MType == metrics.Gauge {
		currentModel.Value = newModel.Value
		return currentModel, nil
	}
	if newModel.Delta != nil && currentModel.Delta != nil {
		newVal := *currentModel.Delta + *newModel.Delta
		currentModel.Delta = &newVal
		return currentModel, nil
	}
	currentModel.Delta = newModel.Delta
	return currentModel, nil
}
