package mocks

import (
	"io"
	"net/http"
)

// MockClient Mock клиент для тестирования
type MockClient struct {
	PostFunc func(url, contentType string, body io.Reader) (*http.Response, error)
	DoFunc   func(req *http.Request) (*http.Response, error)
}

// Post отправляет запрос методом Post
func (m *MockClient) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	return m.PostFunc(url, contentType, body)
}

// Do отправляет запрос
func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}
