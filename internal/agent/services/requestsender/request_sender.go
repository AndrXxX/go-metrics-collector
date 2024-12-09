package requestsender

import (
	"bytes"
	"net/http"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/requestsender/dto"
)

// RequestSender сервис для отправки запросов
type RequestSender struct {
	c    client
	opts []Option
}

// New возвращает сервис RequestSender для отправки запросов
func New(c client, opts ...Option) *RequestSender {
	return &RequestSender{c, opts}
}

// Post отправляет запрос методом Post
func (s *RequestSender) Post(url string, contentType string, data []byte) error {
	params := dto.ParamsDto{
		Buf:     bytes.NewBuffer(data),
		Headers: map[string]string{"Content-Type": contentType},
	}
	for _, opt := range s.opts {
		err := opt(&params)
		if err != nil {
			return err
		}
	}
	r, err := http.NewRequest("POST", url, params.Buf)
	if err != nil {
		return err
	}
	for k, v := range params.Headers {
		r.Header.Set(k, v)
	}

	resp, err := s.c.Do(r)
	if err != nil {
		return err
	}
	if resp != nil && resp.Body != nil {
		return resp.Body.Close()
	}
	return nil
}
