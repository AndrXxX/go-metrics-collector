package storageprovider

type storageProvider struct {
	storages map[string]storage
}

func (sp *storageProvider) RegisterStorage(name string, storage storage) {
	sp.storages[name] = storage
}

func (sp *storageProvider) GetStorage(name string) storage {
	return sp.storages[name]
}

func New() *storageProvider {
	return &storageProvider{
		storages: map[string]storage{},
	}
}
