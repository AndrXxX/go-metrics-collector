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

func (s *storage[T]) Insert(name string, value T) {
	s.store[name] = value
}

func (s *storage[T]) Get(name string) (value T, ok bool) {
	val, found := s.store[name]
	return val, found
}

func (s *storage[T]) All(_ context.Context) map[string]T {
	return s.store
}
