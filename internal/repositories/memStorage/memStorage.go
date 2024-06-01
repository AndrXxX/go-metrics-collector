package memStorage

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func New() MemStorage {
	return MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

func (s *MemStorage) Gauge(metric string, value float64) {
	s.gauge[metric] = value
}

func (s *MemStorage) Counter(metric string, value int64) {
	s.counter[metric] += value
}
