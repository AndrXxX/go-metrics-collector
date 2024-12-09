package filestorage

import "context"

type storageSaver interface {
	Save(ctx context.Context) error
	Restore(ctx context.Context) error
}

type Option func(fs *fileStorage)
