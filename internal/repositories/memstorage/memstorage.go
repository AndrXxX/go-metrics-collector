package memstorage

import (
	"errors"
	"fmt"
)

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

func (s *MemStorage) SetGauge(metric string, value float64) {
	s.gauge[metric] = value
}

func (s *MemStorage) GetGauge(metric string) (value float64, err error) {
	if val, ok := s.gauge[metric]; ok {
		return val, nil
	}
	return 0, errors.New(fmt.Sprintf("value %s not exists", metric))
}

func (s *MemStorage) SetCounter(metric string, value int64) {
	s.counter[metric] += value
}

func (s *MemStorage) GetCounter(metric string) (value int64, err error) {
	if val, ok := s.counter[metric]; ok {
		return val, nil
	}
	return 0, errors.New(fmt.Sprintf("value %s not exists", metric))
}
