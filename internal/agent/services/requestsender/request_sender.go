package requestsender

import (
	"bytes"
	"compress/gzip"
	"fmt"
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

func (s *RequestSender) Post(url string, contentType string, data []byte) error {
	buf, err := s.compress(data)
	if err != nil {
		return err
	}
	r, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return err
	}
	r.Header.Set("Content-Type", contentType)
	r.Header.Set("Content-Encoding", "gzip")
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

func (s *RequestSender) compress(data []byte) (*bytes.Buffer, error) {
	var b bytes.Buffer
	if data == nil {
		return &b, nil
	}
	w := gzip.NewWriter(&b)
	_, err := w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
	}

	return &b, nil
}
