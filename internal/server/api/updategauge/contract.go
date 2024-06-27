package updategauge

type guStorage interface {
	Insert(metric string, value float64)
}
