package configfileparser

type pathProvider interface {
	Fetch() (string, error)
}
