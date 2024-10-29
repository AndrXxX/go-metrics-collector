package dto

// MetricsDto для хранения JSONMetrics
type MetricsDto struct {
	list map[string]JSONMetrics
}

// Get получение метрики JSONMetrics
func (dto *MetricsDto) Get(name string) (JSONMetrics, bool) {
	v, ok := dto.list[name]
	return v, ok
}

// Set установка метрики JSONMetrics
func (dto *MetricsDto) Set(m JSONMetrics) {
	dto.list[m.ID] = m
}

// All получение всех JSONMetrics
func (dto *MetricsDto) All() map[string]JSONMetrics {
	return dto.list
}

// NewMetricsDto возвращает новый экземпляр MetricsDto
func NewMetricsDto() *MetricsDto {
	return &MetricsDto{
		list: map[string]JSONMetrics{},
	}
}
