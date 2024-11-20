package configfile

type pathProvider interface {
	Fetch() (string, error)
}
