package fetchmetrics

type mfStorage[T any] interface {
	All() map[string]T
}
