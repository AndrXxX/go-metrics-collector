package savestoragetask

import "context"

type storageSaver interface {
	Save(ctx context.Context) error
}
