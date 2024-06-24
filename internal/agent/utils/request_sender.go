package utils

import (
	"io"
	"net/http"
)

type Client interface {
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
}

type RequestSender struct {
	ub URLBuilder
	c  Client
}

func NewRequestSender(ub URLBuilder, c Client) *RequestSender {
	return &RequestSender{ub: ub, c: c}
}

func (s *RequestSender) Post(params URLParams, contentType string) error {
	url := s.ub.BuildURL(params)
	resp, err := s.c.Post(url, contentType, nil)
	if err != nil {
		return err
	}
	if resp != nil && resp.Body != nil {
		return resp.Body.Close()
	}
	return nil
}
