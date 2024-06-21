package repositories

type Repository interface {
	SetGauge(metric string, value float64)
	GetGauge(metric string) (value float64, err error)
	GetGaugeAll() map[string]float64
	SetCounter(metric string, value int64)
	GetCounter(metric string) (value int64, err error)
	GetCounterAll() map[string]int64
}
