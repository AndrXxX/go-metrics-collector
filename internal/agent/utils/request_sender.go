package utils

import "net/http"

type RequestSender struct {
	ub URLBuilder
}

func NewRequestSender(ub URLBuilder) *RequestSender {
	return &RequestSender{ub: ub}
}

func (s *RequestSender) Post(params URLParams, contentType string) error {
	url := s.ub.BuildURL(params)
	resp, err := http.Post(url, contentType, nil)
	if err != nil {
		return err
	}
	if resp != nil && resp.Body != nil {
		return resp.Body.Close()
	}
	return nil
}
