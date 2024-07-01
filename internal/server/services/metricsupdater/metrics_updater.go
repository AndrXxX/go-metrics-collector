package metricsupdater

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/metricsvaluesetter"
)

type metricsUpdater struct {
	sp storageProvider
}

func New(sp storageProvider) *metricsUpdater {
	return &metricsUpdater{sp}
}

func (u *metricsUpdater) Update(m *models.Metrics, value string) error {
	storage := u.sp.GetStorage(m.MType)
	if storage == nil {
		return fmt.Errorf("not found storage for metrics type %s", m.MType)
	}
	current, exist := storage.Get(m.ID)
	if !exist {
		current = m
		storage.Insert(current.ID, current)
	}
	setter := metricsvaluesetter.Factory().SetterByType(current.MType)
	return setter.Set(current, value)
}
