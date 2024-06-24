package repositories

type GaugeStorage interface {
	SetGauge(metric string, value float64)
	GetGauge(metric string) (value float64, err error)
	GetGaugeAll() map[string]float64
}

type CounterStorage interface {
	SetCounter(metric string, value int64)
	GetCounter(metric string) (value int64, err error)
	GetCounterAll() map[string]int64
}
