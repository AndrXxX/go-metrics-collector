package metricsvaluesetter

import (
	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
)

type setter interface {
	Set(m *models.Metrics, value string) error
}

type factory struct {
	setters map[string]setter
}

func (f *factory) CounterValueSetter() *counterValueSetter {
	return &counterValueSetter{}
}

func (f *factory) GaugeValueSetter() *gaugeValueSetter {
	return &gaugeValueSetter{}
}

func (f *factory) SetterByType(mType string) (setter, bool) {
	s, ok := f.setters[mType]
	return s, ok
}

func Factory() *factory {
	return &factory{
		setters: map[string]setter{
			metrics.Counter: &counterValueSetter{},
			metrics.Gauge:   &gaugeValueSetter{},
		},
	}
}
