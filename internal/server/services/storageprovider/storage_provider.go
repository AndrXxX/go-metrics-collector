package storageprovider

type storageProvider[T any] struct {
	storages map[string]T
}

func (sp *storageProvider[T]) RegisterStorage(name string, storage T) {
	sp.storages[name] = storage
}

func (sp *storageProvider[T]) GetStorage(name string) T {
	return sp.storages[name]
}

func New[T any]() *storageProvider[T] {
	return &storageProvider[T]{
		storages: map[string]T{},
	}
}
