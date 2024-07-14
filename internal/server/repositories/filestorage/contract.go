package filestorage

type storageSaver interface {
	Save() error
	Restore() error
}
