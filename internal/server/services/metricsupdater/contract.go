package metricsupdater

import "context"

type storage[T any] interface {
	Insert(ctx context.Context, metric string, value T)
	Get(ctx context.Context, metric string) (value T, ok bool)
}
