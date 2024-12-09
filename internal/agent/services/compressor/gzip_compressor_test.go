package compressor

import (
	"bytes"
	"fmt"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

type testBuffer struct {
	err             error
	errOnWriteBytes []byte
}

func (b *testBuffer) Read(p []byte) (n int, err error) {
	return len(p), b.err
}

func (b *testBuffer) Write(p []byte) (n int, err error) {
	if b.errOnWriteBytes != nil && slices.Equal(p, b.errOnWriteBytes) {
		return 0, fmt.Errorf("expected error writing to buffer")
	}
	return len(p), b.err
}

func TestGzipCompressor_Compress(t *testing.T) {
	tests := []struct {
		name    string
		buffer  buffer
		data    []byte
		wantErr bool
	}{
		{
			name:    "Test with empty data",
			buffer:  &testBuffer{},
			data:    nil,
			wantErr: false,
		},
		{
			name:    "Test with error on write data",
			buffer:  &testBuffer{err: fmt.Errorf("test error")},
			data:    []byte("test"),
			wantErr: true,
		},
		{
			name:    "Test with error on close",
			buffer:  &testBuffer{errOnWriteBytes: []byte{0, 0, 255, 255}},
			data:    []byte("test"),
			wantErr: true,
		},
		{
			name:    "Test with data",
			buffer:  bytes.NewBuffer(nil),
			data:    []byte("test"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := GzipCompressor{Buff: tt.buffer}
			_, err := c.Compress(tt.data)
			require.Equal(t, tt.wantErr, err != nil)
		})
	}
}
