package requestsender

import (
	"net/http"
)

type client interface {
	Do(req *http.Request) (*http.Response, error)
}

type hashGenerator interface {
	Generate(key string, data []byte) string
}
