package dbchecker

import "context"

type conn interface {
	PingContext(ctx context.Context) error
}
