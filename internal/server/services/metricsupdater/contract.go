package metricsupdater

type storage[T any] interface {
	Insert(metric string, value T)
	Get(metric string) (value T, ok bool)
}
