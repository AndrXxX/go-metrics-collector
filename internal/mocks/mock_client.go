package mocks

import (
	"io"
	"net/http"
)

type MockClient struct {
	PostFunc func(url, contentType string, body io.Reader) (*http.Response, error)
}

func (m *MockClient) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	return m.PostFunc(url, contentType, body)
}
