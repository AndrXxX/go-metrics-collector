package fetchmetrics

type storage[T any] interface {
	All() map[string]T
}
