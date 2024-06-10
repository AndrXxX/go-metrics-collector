package metrics

type Metrics struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

func NewMetrics() *Metrics {
	return &Metrics{
		Gauge:   map[string]float64{},
		Counter: map[string]int64{},
	}
}
