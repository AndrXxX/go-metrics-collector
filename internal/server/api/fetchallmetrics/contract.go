package fetchallmetrics

import "context"

type storage[T any] interface {
	All(ctx context.Context) map[string]T
}
