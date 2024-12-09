package dealoptions

type hashGenerator interface {
	Generate(key string, data []byte) string
}
