package compressor

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

type GzipCompressor struct {
}

func (c GzipCompressor) Compress(data []byte) (io.Reader, error) {
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
