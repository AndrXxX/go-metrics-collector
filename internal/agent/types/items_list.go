package types

type ItemsList[T any] []T

func (list *ItemsList[T]) Add(v T) {
	*list = append(*list, v)
}
