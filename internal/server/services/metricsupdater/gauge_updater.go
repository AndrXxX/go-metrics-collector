package metricsupdater

import (
	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"github.com/AndrXxX/go-metrics-collector/internal/server/repositories"
	"strconv"
)

type gaugeUpdater struct {
	storage repositories.Storage[*models.Metrics]
}

func NewGaugeUpdater(storage repositories.Storage[*models.Metrics]) *gaugeUpdater {
	return &gaugeUpdater{storage: storage}
}

func (u *gaugeUpdater) Update(name string, value string) error {
	converted, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	u.storage.Insert(name, &models.Metrics{
		ID:    name,
		MType: metrics.Counter,
		Value: &converted,
	})
	return nil
}
