package memory

import "context"

type storage[T any] struct {
	store map[string]T
}

func New[T any]() storage[T] {
	return storage[T]{
		store: map[string]T{},
	}
}

func (s *storage[T]) Insert(_ context.Context, name string, value T) {
	s.store[name] = value
}

func (s *storage[T]) Get(_ context.Context, name string) (value T, ok bool) {
	val, found := s.store[name]
	return val, found
}

func (s *storage[T]) All(_ context.Context) map[string]T {
	return s.store
}

func (s *storage[T]) Delete(_ context.Context, name string) (ok bool) {
	delete(s.store, name)
	return true
}
