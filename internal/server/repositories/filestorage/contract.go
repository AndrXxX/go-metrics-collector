package filestorage

import "context"

type storageSaver interface {
	Save(ctx context.Context) error
	Restore() error
}
