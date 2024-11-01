package metricsuploader

import "github.com/AndrXxX/go-metrics-collector/internal/agent/types"

type urlBuilder interface {
	Build(params types.URLParams) string
}

type requestSender interface {
	Post(url string, contentType string, data []byte) error
}
