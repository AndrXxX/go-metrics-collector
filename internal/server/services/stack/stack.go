package stack

type stack[T any] struct {
	elements []T
}

func (s *stack[T]) Push(value T) {
	s.elements = append(s.elements, value)
}

func (s *stack[T]) Pop() (T, bool) {
	if len(s.elements) == 0 {
		var zero T
		return zero, false
	}
	last := s.elements[len(s.elements)-1]
	s.elements = s.elements[:len(s.elements)-1]
	return last, true
}

func NewFromSlice[T any](elements []T) *stack[T] {
	return &stack[T]{elements: elements}
}
