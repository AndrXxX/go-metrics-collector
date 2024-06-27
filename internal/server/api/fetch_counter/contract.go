package fetch_counter

type cfStorage interface {
	Get(metric string) (value int64, ok bool)
}
