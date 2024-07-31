package stack

import "sync"

type Stack[T any] struct {
	elements []T
	m        sync.Mutex
}

func (s *Stack[T]) Push(value T) {
	s.m.Lock()
	defer s.m.Unlock()
	s.elements = append(s.elements, value)
}

func (s *Stack[T]) Pop() (T, bool) {
	s.m.Lock()
	defer s.m.Unlock()
	if len(s.elements) == 0 {
		var zero T
		return zero, false
	}
	last := s.elements[len(s.elements)-1]
	s.elements = s.elements[:len(s.elements)-1]
	return last, true
}

func (s *Stack[T]) Shift() (T, bool) {
	s.m.Lock()
	defer s.m.Unlock()
	if len(s.elements) == 0 {
		var zero T
		return zero, false
	}
	first := s.elements[0]
	s.elements = s.elements[1:len(s.elements)]
	return first, true
}

func (s *Stack[T]) All() []T {
	s.m.Lock()
	defer s.m.Unlock()
	return s.elements
}

func (s *Stack[T]) Copy() *Stack[T] {
	s.m.Lock()
	defer s.m.Unlock()
	return NewFromSlice(s.elements)
}

func NewFromSlice[T any](elements []T) *Stack[T] {
	return &Stack[T]{elements: elements}
}
