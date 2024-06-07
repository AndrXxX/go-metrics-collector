package repositories

type Repository interface {
	SetGauge(metric string, value float64)
	SetCounter(metric string, value int64)
}
