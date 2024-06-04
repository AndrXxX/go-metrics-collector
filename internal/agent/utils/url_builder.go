package utils

type URLBuilder interface {
	BuildURL(params URLParams) string
}

type URLParams map[string]any
