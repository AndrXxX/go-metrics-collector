package requestsender

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// RequestSender сервис для отправки запросов
type RequestSender struct {
	c    client
	hg   hashGenerator
	comp dataCompressor
	key  string
}

// New возвращает сервис RequestSender для отправки запросов
func New(c client, hg hashGenerator, key string, comp dataCompressor) *RequestSender {
	return &RequestSender{c, hg, comp, key}
}

// Post отправляет запрос методом Post
func (s *RequestSender) Post(url string, contentType string, data []byte) error {
	buf, err := s.comp.Compress(data)
	if err != nil {
		return err
	}
	var encoded []byte
	if s.key != "" {
		encoded, err = io.ReadAll(buf)
		if err != nil {
			return fmt.Errorf("error on read encoded data: %w", err)
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
