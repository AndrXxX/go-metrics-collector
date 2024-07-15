package repositories

import "context"

type Storage[T any] interface {
	Insert(metric string, value T)
	Get(metric string) (value T, ok bool)
	All(ctx context.Context) map[string]T
}

type StorageShutdowner interface {
	Shutdown(ctx context.Context) error
}

type StorageSaver interface {
	Save(ctx context.Context) error
}
