package fetchcounter

type cfStorage interface {
	Get(metric string) (value int64, ok bool)
}
