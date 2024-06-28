package dto

type MetricsDto struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

func NewMetricsDto() *MetricsDto {
	return &MetricsDto{
		Gauge:   map[string]float64{},
		Counter: map[string]int64{},
	}
}
