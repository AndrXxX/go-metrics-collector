package dto

type MetricsDto struct {
	Gauge   map[string]float64
	Counter map[string]int64
	list    map[string]*JSONMetrics
}

func (dto *MetricsDto) Get(name string) (*JSONMetrics, bool) {
	v, ok := dto.list[name]
	return v, ok
}

func (dto *MetricsDto) Set(m *JSONMetrics) {
	dto.list[m.ID] = m
}

func (dto *MetricsDto) All() map[string]*JSONMetrics {
	return dto.list
}

func NewMetricsDto() *MetricsDto {
	return &MetricsDto{
		Gauge:   map[string]float64{},
		Counter: map[string]int64{},
		list:    map[string]*JSONMetrics{},
	}
}
