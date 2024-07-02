package storagesaver

type storage[T any] interface {
	All() map[string]T
	Insert(metric string, value T)
}
