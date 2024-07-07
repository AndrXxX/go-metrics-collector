package savestoragetask

type storageSaver interface {
	Save() error
}
