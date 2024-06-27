package update_gauge

type guStorage interface {
	Insert(metric string, value float64)
}
