package middlewares

type SHA256hashGenerator interface {
	Generate(data []byte) string
}
