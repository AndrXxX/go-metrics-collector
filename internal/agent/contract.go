package agent

type Option func(a *agent)

type hashGenerator interface {
	Generate(key string, data []byte) string
}
