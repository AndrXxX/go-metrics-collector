package configpath

type fetcher interface {
	Fetch() (string, error)
}
