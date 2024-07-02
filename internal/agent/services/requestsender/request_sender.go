package requestsender

import (
	"compress/gzip"
	"io"
	"net/http"
)

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

type RequestSender struct {
	c Client
}

func New(c Client) *RequestSender {
	return &RequestSender{c}
}

func (s *RequestSender) Post(url string, contentType string, body io.Reader) error {
	var err error
	if body != nil {
		body, err = gzip.NewReader(body)
	}
	if err != nil {
		return err
	}
	r, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	r.Header.Set("Content-Type", contentType)
	r.Header.Set("Accept-Encoding", "gzip")
	resp, err := s.c.Do(r)
	if err != nil {
		return err
	}
	if resp != nil && resp.Body != nil {
		return resp.Body.Close()
	}
	return nil
}
