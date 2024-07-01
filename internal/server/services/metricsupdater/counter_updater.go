package metricsupdater

import (
	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"github.com/AndrXxX/go-metrics-collector/internal/server/repositories"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/metricsvaluesetter"
)

type counterUpdater struct {
	storage repositories.Storage[*models.Metrics]
	setter  metricsSetter
}

func NewCounterUpdater(storage repositories.Storage[*models.Metrics]) *counterUpdater {
	return &counterUpdater{storage: storage, setter: &metricsvaluesetter.CounterValueSetter{}}
}

func (u *counterUpdater) Update(name string, value string) error {
	current, exist := u.storage.Get(name)
	if !exist {
		current = &models.Metrics{
			ID:    name,
			MType: metrics.Counter,
		}
		u.storage.Insert(name, current)
	}
	return u.setter.Set(current, value)
}
