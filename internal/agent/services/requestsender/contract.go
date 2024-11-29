package requestsender

import (
	"io"
	"net/http"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/requestsender/dto"
)

type client interface {
	Do(req *http.Request) (*http.Response, error)
}

type hashGenerator interface {
	Generate(key string, data []byte) string
}

type dataCompressor interface {
	Compress(in []byte) (io.Reader, error)
}

type Option func(*dto.ParamsDto) error
