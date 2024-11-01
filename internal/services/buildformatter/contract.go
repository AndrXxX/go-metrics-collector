package buildformatter

type buffer interface {
	Write(p []byte) (n int, err error)
	String() string
}
