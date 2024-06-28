package metricscollector

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/dto"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
	"math/rand"
	"runtime"
)

type metricsCollector struct {
	ml *config.MetricsList
}

func (c *metricsCollector) Execute(result dto.MetricsDto) error {
	memStatsDto := dto.NewMemStatsDto()
	runtime.ReadMemStats(&memStatsDto.Stats)

	for _, name := range *c.ml {
		val, err := memStatsDto.GetValue(name)
		if err != nil {
			logger.Log.Error(fmt.Sprintf("Failed to get value for metric %s: %s", name, err.Error()))
			continue
		}
		result.Gauge[name] = val
	}
	result.Counter[metrics.PollCount]++
	result.Gauge[metrics.RandomValue] = rand.Float64()
	return nil
}

func New(ml *config.MetricsList) *metricsCollector {
	return &metricsCollector{ml}
}
