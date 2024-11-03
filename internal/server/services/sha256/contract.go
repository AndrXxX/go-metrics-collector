package sha256

type hashGenerator interface {
	Generate(key string, data []byte) string
}
