package dbping

import "context"

type dbChecker interface {
	Check(context.Context) error
}
