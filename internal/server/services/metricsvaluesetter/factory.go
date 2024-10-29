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

// CounterValueSetter возвращает сервис для установки значения метрики типа counter
func (f *factory) CounterValueSetter() *counterValueSetter {
	return &counterValueSetter{}
}

// GaugeValueSetter возвращает сервис для установки значения метрики типа gauge
func (f *factory) GaugeValueSetter() *gaugeValueSetter {
	return &gaugeValueSetter{}
}

// SetterByType возвращает сервис для установки значения метрики соответствующего типа
func (f *factory) SetterByType(mType string) (setter, bool) {
	s, ok := f.setters[mType]
	return s, ok
}

// Factory возвращает фабрику сервисов для установки значений метрик
func Factory() *factory {
	return &factory{
		setters: map[string]setter{
			metrics.Counter: &counterValueSetter{},
			metrics.Gauge:   &gaugeValueSetter{},
		},
	}
}
