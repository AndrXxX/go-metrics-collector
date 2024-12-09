package gzipcompressor

import (
	"compress/gzip"
	"io"
)

type compressReader struct {
	r  io.ReadCloser
	zr io.ReadCloser
}

// NewCompressReader возвращает сервис для чтения сжатых данных gzip
func NewCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

// Read имплементация io.Reader
func (c *compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close имплементация io.Closer
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
