package middlewares

type SHA256hashGenerator interface {
	Generate(key string, data []byte) string
}
