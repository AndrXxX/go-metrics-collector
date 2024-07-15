package storagesaver

import "context"

type storage[T any] interface {
	All(ctx context.Context) map[string]T
	Insert(metric string, value T)
}
