package metricsupdater

import (
	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"github.com/AndrXxX/go-metrics-collector/internal/server/repositories"
	"strconv"
)

type counterUpdater struct {
	storage repositories.Storage[*models.Metrics]
}

func NewCounterUpdater(storage repositories.Storage[*models.Metrics]) *counterUpdater {
	return &counterUpdater{storage: storage}
}

func (u *counterUpdater) Update(name string, value string) error {
	converted, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	current, exist := u.storage.Get(name)
	if !exist {
		u.storage.Insert(name, &models.Metrics{
			ID:    name,
			MType: metrics.Counter,
			Delta: &converted,
		})
		return nil
	}
	current.Delta = &converted
	return nil
}
