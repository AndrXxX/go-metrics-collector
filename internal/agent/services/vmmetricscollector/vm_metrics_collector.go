package vmmetricscollector

import (
	"fmt"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/dto"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
)

type collector struct {
}

func (c *collector) Execute(result dto.MetricsDto) error {
	v, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	tm := float64(v.Total)
	result.Set(dto.JSONMetrics{ID: metrics.TotalMemory, MType: metrics.Gauge, Value: &tm})

	fm := float64(v.Free)
	result.Set(dto.JSONMetrics{ID: metrics.FreeMemory, MType: metrics.Gauge, Value: &fm})

	cpuList, _ := cpu.Percent(0, true)
	for i, percent := range cpuList {
		id := fmt.Sprintf("%s%d", metrics.CPUutilization, i+1)
		result.Set(dto.JSONMetrics{ID: id, MType: metrics.Gauge, Value: &percent})
	}
	return nil
}

func (c *collector) Collect(results chan<- dto.MetricsDto) error {
	m := dto.NewMetricsDto()
	err := c.Execute(*m)
	if err != nil {
		return err
	}
	results <- *m
	close(results)
	return nil
}

func New() *collector {
	return &collector{}
}
