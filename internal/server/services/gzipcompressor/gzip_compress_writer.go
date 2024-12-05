package gzipcompressor

import (
	"compress/gzip"
	"io"
	"net/http"
)

type compressWriter struct {
	w  http.ResponseWriter
	zw io.WriteCloser
}

// NewCompressWriter возвращает сервис для записи сжатых данных gzip
func NewCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// Header имплементация http.ResponseWriter
func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

// Write имплементация http.ResponseWriter
func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

// WriteHeader имплементация http.ResponseWriter
func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close имплементация io.Closer
func (c *compressWriter) Close() error {
	return c.zw.Close()
}
