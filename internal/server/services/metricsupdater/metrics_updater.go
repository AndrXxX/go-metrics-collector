package metricsupdater

import "github.com/AndrXxX/go-metrics-collector/internal/server/models"

type metricsUpdater struct {
	storage storage
	setter  metricsSetter
	mType   string
}

func New(storage storage, setter metricsSetter, mType string) *metricsUpdater {
	return &metricsUpdater{storage, setter, mType}
}

func (u *metricsUpdater) Update(name string, value string) error {
	current, exist := u.storage.Get(name)
	if !exist {
		current = &models.Metrics{
			ID:    name,
			MType: u.mType,
		}
		u.storage.Insert(name, current)
	}
	return u.setter.Set(current, value)
}
