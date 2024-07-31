package memory

import (
	"context"
	"sync"
)

type storage[T any] struct {
	m sync.Map
}

func New[T any]() storage[T] {
	return storage[T]{}
}

func (s *storage[T]) Insert(_ context.Context, name string, value T) {
	s.m.Store(name, value)
}

func (s *storage[T]) Get(_ context.Context, name string) (value T, ok bool) {
	val, found := s.m.Load(name)
	if val, ok := val.(T); ok {
		return val, ok
	}
	var zero T
	return zero, found
}

func (s *storage[T]) All(_ context.Context) map[string]T {
	result := make(map[string]T)
	s.m.Range(func(k, v interface{}) bool {
		result[k.(string)] = v.(T)
		return true
	})
	return result
}

func (s *storage[T]) Delete(_ context.Context, name string) (ok bool) {
	s.m.Delete(name)
	return true
}
