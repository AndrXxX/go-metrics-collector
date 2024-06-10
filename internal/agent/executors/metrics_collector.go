package executors

import (
	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/metrics"
	me "github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"math/rand"
	"reflect"
	"runtime"
)

type metricsCollector struct {
	ml *config.MetricsList
}

func (c *metricsCollector) Execute(result metrics.Metrics) error {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)
	values := reflect.ValueOf(stats)
	for _, name := range *c.ml {
		field := values.FieldByName(name)
		if field.Kind() == reflect.Ptr {
			field = field.Elem()
		}
		if field.Kind() == reflect.Invalid {
			continue
		}
		if field.CanFloat() {
			result.Gauge[name] = field.Float()
			continue
		}
		if field.CanInt() {
			result.Gauge[name] = float64(field.Int())
			continue
		}
		if field.CanUint() {
			result.Gauge[name] = float64(field.Uint())
			continue
		}
	}
	result.Counter[me.PollCount]++
	result.Gauge[me.RandomValue] = rand.Float64()
	return nil
}

func NewCollector(ml *config.MetricsList) Executors {
	return &metricsCollector{
		ml: ml,
	}
}
