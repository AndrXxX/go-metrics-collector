package types

type ItemsList[T any] []T

func (list *ItemsList[T]) Add(p T) {
	*list = append(*list, p)
}
