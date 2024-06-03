package executors

import (
	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/metrics"
	"math/rand"
	rm "runtime/metrics"
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

func (c *metricsCollector) buildSamples(list []string) []rm.Sample {
	result := make([]rm.Sample, 0)
	for _, name := range list {
		result = append(result, rm.Sample{Name: name})
	}
	return result
}

func (c *metricsCollector) Execute(result metrics.Metrics) error {
	if !c.canExecute() {
		return nil
	}
	samples := c.buildSamples(c.config.Metrics)
	rm.Read(samples)
	result.Counter["PollCount"]++
	result.Gauge["RandomValue"] = rand.Float64()
	for _, sample := range samples {
		if sample.Value.Kind() == rm.KindFloat64 {
			result.Gauge[sample.Name] = sample.Value.Float64()
		} else if sample.Value.Kind() == rm.KindUint64 {
			result.Gauge[sample.Name] = float64(sample.Value.Uint64())
		}
	}
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
