package updatecounter

type updater interface {
	Update(name string, value int64)
}
