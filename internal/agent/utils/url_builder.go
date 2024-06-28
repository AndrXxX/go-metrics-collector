package utils

import "github.com/AndrXxX/go-metrics-collector/internal/agent/types"

type URLBuilder interface {
	Build(params types.URLParams) string
}
