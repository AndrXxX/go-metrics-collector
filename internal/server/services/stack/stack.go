package stack

type Stack[T any] struct {
	elements []T
}

func (s *Stack[T]) Push(value T) {
	s.elements = append(s.elements, value)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.elements) == 0 {
		var zero T
		return zero, false
	}
	last := s.elements[len(s.elements)-1]
	s.elements = s.elements[:len(s.elements)-1]
	return last, true
}

func (s *Stack[T]) All() []T {
	return s.elements
}

func (s *Stack[T]) Copy() *Stack[T] {
	return NewFromSlice(s.elements)
}

func NewFromSlice[T any](elements []T) *Stack[T] {
	return &Stack[T]{elements: elements}
}
