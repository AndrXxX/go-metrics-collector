package counter

import "github.com/AndrXxX/go-metrics-collector/internal/server/repositories"

type updater struct {
	storage repositories.Storage[int64]
}

func New(storage repositories.Storage[int64]) updater {
	return updater{storage: storage}
}

func (u *updater) Update(name string, value int64) {
	current, ok := u.storage.Get(name)
	if !ok {
		current = 0
	}
	u.storage.Insert(name, current+value)
}
