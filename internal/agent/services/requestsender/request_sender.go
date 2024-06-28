package requestsender

import (
	"github.com/AndrXxX/go-metrics-collector/internal/agent/types"
	"io"
	"net/http"
)

type Client interface {
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
}

type RequestSender struct {
	ub urlBuilder
	c  Client
}

func New(ub urlBuilder, c Client) *RequestSender {
	return &RequestSender{ub: ub, c: c}
}

func (s *RequestSender) Post(params types.URLParams, contentType string) error {
	url := s.ub.Build(params)
	resp, err := s.c.Post(url, contentType, nil)
	if err != nil {
		return err
	}
	if resp != nil && resp.Body != nil {
		return resp.Body.Close()
	}
	return nil
}
