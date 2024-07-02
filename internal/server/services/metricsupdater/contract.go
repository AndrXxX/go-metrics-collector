package metricsupdater

type storageProvider[T any] interface {
	GetStorage(name string) T
}
