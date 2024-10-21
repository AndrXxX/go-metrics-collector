package memory

import (
	"context"
	"sync"
)

type storage[T any] struct {
	m sync.Map
}

// New возвращает хранилище метрик в оперативной памяти
func New[T any]() storage[T] {
	return storage[T]{}
}

// Insert вставляет запись
func (s *storage[T]) Insert(_ context.Context, name string, value T) {
	s.m.Store(name, value)
}

// Get извлекает запись
func (s *storage[T]) Get(_ context.Context, name string) (value T, ok bool) {
	val, found := s.m.Load(name)
	if val, ok := val.(T); ok {
		return val, ok
	}
	var zero T
	return zero, found
}

// All извлекает все записи
func (s *storage[T]) All(_ context.Context) map[string]T {
	result := make(map[string]T)
	s.m.Range(func(k, v interface{}) bool {
		result[k.(string)] = v.(T)
		return true
	})
	return result
}

// Delete удаляет запись
func (s *storage[T]) Delete(_ context.Context, name string) (ok bool) {
	s.m.Delete(name)
	return true
}
