package requestsender

import "github.com/AndrXxX/go-metrics-collector/internal/agent/types"

type urlBuilder interface {
	Build(params types.URLParams) string
}
