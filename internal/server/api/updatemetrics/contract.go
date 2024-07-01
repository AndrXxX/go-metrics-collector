package updatemetrics

type updater interface {
	Update(name string, value string) error
}
