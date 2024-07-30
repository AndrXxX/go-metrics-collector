package memory

import (
	"context"
	"sync"
)

type storage[T any] struct {
	store map[string]T
	m     sync.Mutex
}

func New[T any]() storage[T] {
	return storage[T]{
		store: map[string]T{},
	}
}

func (s *storage[T]) Insert(_ context.Context, name string, value T) {
	s.m.Lock()
	s.store[name] = value
	s.m.Unlock()
}

func (s *storage[T]) Get(_ context.Context, name string) (value T, ok bool) {
	s.m.Lock()
	val, found := s.store[name]
	s.m.Unlock()
	return val, found
}

func (s *storage[T]) All(_ context.Context) map[string]T {
	return s.store
}

func (s *storage[T]) Delete(_ context.Context, name string) (ok bool) {
	s.m.Lock()
	delete(s.store, name)
	s.m.Unlock()
	return true
}
