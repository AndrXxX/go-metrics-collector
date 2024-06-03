package executors

import (
	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/metrics"
	"math/rand"
	"reflect"
	"runtime"
	"slices"
	"time"
)

type metricsCollector struct {
	interval     time.Duration
	lastExecuted time.Time
	config       *config.Config
}

func (c *metricsCollector) canExecute() bool {
	return time.Since(c.lastExecuted) >= c.interval
}

func (c *metricsCollector) Execute(result metrics.Metrics) error {
	if !c.canExecute() {
		return nil
	}
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)
	values := reflect.ValueOf(stats)
	types := values.Type()
	for i := 0; i < values.NumField(); i++ {
		if !slices.Contains(c.config.Metrics, types.Field(i).Name) {
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
	result.Counter["PollCount"]++
	result.Gauge["RandomValue"] = rand.Float64()
	c.lastExecuted = time.Now()
	return nil
}

func NewCollector(interval time.Duration, config *config.Config) Executors {
	return &metricsCollector{
		interval:     interval,
		lastExecuted: time.Now(),
		config:       config,
	}
}
