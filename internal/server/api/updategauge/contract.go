package updategauge

type updater interface {
	Update(name string, value string) error
}
