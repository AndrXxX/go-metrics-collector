package requestsender

import (
	"io"
	"net/http"
)

type Client interface {
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
}

type RequestSender struct {
	c Client
}

func New(c Client) *RequestSender {
	return &RequestSender{c}
}

func (s *RequestSender) Post(url string, contentType string) error {
	resp, err := s.c.Post(url, contentType, nil)
	if err != nil {
		return err
	}
	if resp != nil && resp.Body != nil {
		return resp.Body.Close()
	}
	return nil
}
