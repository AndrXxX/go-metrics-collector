package metricsupdater

import (
	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"github.com/AndrXxX/go-metrics-collector/internal/server/repositories"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/metricsvaluesetter"
)

type gaugeUpdater struct {
	storage repositories.Storage[*models.Metrics]
	setter  metricsSetter
}

func NewGaugeUpdater(storage repositories.Storage[*models.Metrics]) *gaugeUpdater {
	return &gaugeUpdater{storage: storage, setter: metricsvaluesetter.Factory().GaugeValueSetter()}
}

func (u *gaugeUpdater) Update(name string, value string) error {
	current, exist := u.storage.Get(name)
	if !exist {
		current = &models.Metrics{
			ID:    name,
			MType: metrics.Gauge,
		}
		u.storage.Insert(name, current)
	}
	return u.setter.Set(current, value)
}
