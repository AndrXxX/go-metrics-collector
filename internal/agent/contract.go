package agent

import (
	"crypto/tls"
	"net/http"
)

type Option func(a *agent)

type hashGenerator interface {
	Generate(key string, data []byte) string
}

type tlsConfigProvider interface {
	Fetch() (*tls.Config, error)
}

type clientProvider interface {
	Fetch() (*http.Client, error)
}
