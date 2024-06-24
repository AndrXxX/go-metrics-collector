package executors

import (
	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/metrics"
	me "github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"math/rand"
	"reflect"
	// TODO: Perederey Использование reflect и slices.Contains в metricsCollector может быть избыточным. Попробуй найти более простой способ обработки полей структуры.
	"runtime"
	"slices"
)

type metricsCollector struct {
	ml *config.MetricsList
}

func (c *metricsCollector) Execute(result metrics.Metrics) error {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)
	values := reflect.ValueOf(stats)
	types := values.Type()
	for i := 0; i < values.NumField(); i++ {
		if !slices.Contains(*c.ml, types.Field(i).Name) {
			continue
		}
		if values.Field(i).CanFloat() {
			result.Gauge[types.Field(i).Name] = values.Field(i).Float()
			continue
		}
		if values.Field(i).CanInt() {
			result.Gauge[types.Field(i).Name] = float64(values.Field(i).Int())
			continue
		}
		if values.Field(i).CanUint() {
			result.Gauge[types.Field(i).Name] = float64(values.Field(i).Uint())
			continue
		}
		println(types.Field(i).Name, values.Field(i).String())
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
