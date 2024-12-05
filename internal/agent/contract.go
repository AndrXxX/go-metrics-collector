package agent

import "crypto/tls"

type Option func(a *agent)

type hashGenerator interface {
	Generate(key string, data []byte) string
}

type tlsConfigProvider interface {
	Fetch() (*tls.Config, error)
}
