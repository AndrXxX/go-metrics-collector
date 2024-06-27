package update_counter

type updater interface {
	Update(name string, value int64)
}
