package metricsvaluesetter

import (
	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
)

type setter interface {
	Set(m *models.Metrics, value string) error
}

type factory struct {
}

func (f *factory) CounterValueSetter() *counterValueSetter {
	return &counterValueSetter{}
}

func (f *factory) GaugeValueSetter() *gaugeValueSetter {
	return &gaugeValueSetter{}
}

func (f *factory) SetterByType(mType string) setter {
	switch mType {
	case metrics.Counter:
		return f.CounterValueSetter()
	case metrics.Gauge:
		return f.GaugeValueSetter()
	default:
		return nil
	}
}

func Factory() *factory {
	return &factory{}
}
