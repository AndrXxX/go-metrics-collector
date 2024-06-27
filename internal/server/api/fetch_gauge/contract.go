package fetch_gauge

type gfStorage interface {
	Get(metric string) (value float64, ok bool)
}
