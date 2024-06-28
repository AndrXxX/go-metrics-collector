package utils

type URLBuilder interface {
	Build(params URLParams) string
}

type URLParams map[string]any
