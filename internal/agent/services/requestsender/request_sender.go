package requestsender

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
)

// RequestSender сервис для отправки запросов
type RequestSender struct {
	c   client
	hg  hashGenerator
	key string
}

// New возвращает сервис RequestSender для отправки запросов
func New(c client, hg hashGenerator, key string) *RequestSender {
	return &RequestSender{c, hg, key}
}

// Post отправляет запрос методом Post
func (s *RequestSender) Post(url string, contentType string, data []byte) error {
	buf, err := s.compress(data)
	if err != nil {
		return err
	}
	var encoded []byte
	if s.key != "" {
		encoded, err = io.ReadAll(buf)
		if err != nil {
			logger.Log.Error("Error on read encoded data", zap.Error(err))
		}
		buf = bytes.NewBuffer(encoded)
	}

	r, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return err
	}
	r.Header.Set("Content-Type", contentType)
	r.Header.Set("Content-Encoding", "gzip")
	r.Header.Set("Accept-Encoding", "gzip")
	if s.key != "" {
		r.Header.Set("HashSHA256", s.hg.Generate(s.key, encoded))
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
