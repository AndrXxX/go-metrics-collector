package dto

type MetricsDto struct {
	list map[string]JSONMetrics
}

func (dto *MetricsDto) Get(name string) (JSONMetrics, bool) {
	v, ok := dto.list[name]
	return v, ok
}

func (dto *MetricsDto) Set(m JSONMetrics) {
	dto.list[m.ID] = m
}

func (dto *MetricsDto) All() map[string]JSONMetrics {
	return dto.list
}

func NewMetricsDto() *MetricsDto {
	return &MetricsDto{
		list: map[string]JSONMetrics{},
	}
}
