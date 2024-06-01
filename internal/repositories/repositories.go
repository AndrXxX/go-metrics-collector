package repositories

type Repository interface {
	Gauge(metric string, value float64)
	Counter(metric string, value int64)
}
