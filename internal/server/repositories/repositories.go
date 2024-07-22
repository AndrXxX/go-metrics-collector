package repositories

import "context"

type Storage[T any] interface {
	Insert(ctx context.Context, metric string, value T)
	Get(ctx context.Context, metric string) (value T, ok bool)
	All(ctx context.Context) map[string]T
	Delete(ctx context.Context, metric string) (ok bool)
}

type StorageShutdowner interface {
	Shutdown(ctx context.Context) error
}

type StorageSaver interface {
	Save(ctx context.Context) error
}
