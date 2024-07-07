package fetchallmetrics

type storage[T any] interface {
	All() map[string]T
}
