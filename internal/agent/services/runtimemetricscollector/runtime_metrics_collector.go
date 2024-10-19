package runtimemetricscollector

import (
	"math/rand"
	"runtime"

	"go.uber.org/zap"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/dto"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
)

type collector struct {
	ml *config.MetricsList
}

func (c *collector) execute(result dto.MetricsDto) {
	ms := runtime.MemStats{}
	runtime.ReadMemStats(&ms)
	memStatsDto := dto.NewMemStatsDto(&ms)

	for _, name := range *c.ml {
		metricFn, ok := memStatsDto.FetchGetter(name)
		if !ok {
			logger.Log.Error("Failed to fetch value getter for metric", zap.String("name", name))
			continue
		}
		v := metricFn()
		result.Set(dto.JSONMetrics{ID: name, MType: metrics.Gauge, Value: &v})
	}
	var pollVal int64
	if curPoll, ok := result.Get(metrics.PollCount); ok {
		pollVal = *curPoll.Delta + 1
	} else {
		pollVal = 1
	}
	result.Set(dto.JSONMetrics{ID: metrics.PollCount, MType: metrics.Counter, Delta: &pollVal})
	randVal := rand.Float64()
	result.Set(dto.JSONMetrics{ID: metrics.RandomValue, MType: metrics.Gauge, Value: &randVal})
}

// Collect собирает runtime метрики и отправляет их в канал results
func (c *collector) Collect(results chan<- dto.MetricsDto) error {
	m := dto.NewMetricsDto()
	c.execute(*m)
	results <- *m
	close(results)
	return nil
}

// New возвращает сервис для сбора runtime метрик
func New(ml *config.MetricsList) *collector {
	return &collector{ml}
}
