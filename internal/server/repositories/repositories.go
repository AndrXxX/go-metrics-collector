package repositories

import "context"

// Storage интерфейс хранилища
type Storage[T any] interface {
	Insert(ctx context.Context, metric string, value T)
	Get(ctx context.Context, metric string) (value T, ok bool)
	All(ctx context.Context) map[string]T
	Delete(ctx context.Context, metric string) (ok bool)
}

// StorageShutdowner интерфейс хранилища с завершением работы
type StorageShutdowner interface {
	Shutdown(ctx context.Context) error
}

// StorageSaver интерфейс хранилища с сохранением данных
type StorageSaver interface {
	Save(ctx context.Context) error
}
