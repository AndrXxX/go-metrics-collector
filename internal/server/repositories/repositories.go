package repositories

type Storage[T any] interface {
	Insert(metric string, value T)
	Get(metric string) (value T, ok bool)
	All() map[string]T
}

type StorageShutdowner interface {
	Shutdown() error
}

type StorageSaver interface {
	Save() error
}
