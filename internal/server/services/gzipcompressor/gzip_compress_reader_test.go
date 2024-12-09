package gzipcompressor

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type closableReadableMock struct {
	closed bool
	err    error
	p      []byte
}

func (m *closableReadableMock) Close() error {
	m.closed = true
	return m.err
}

func (m *closableReadableMock) Read(p []byte) (n int, err error) {
	m.p = p
	return len(p), m.err
}

var gzipData = []byte{
	0x1f, 0x8b, 0x08, 0x08, 0xf7, 0x5e, 0x14, 0x4a,
	0x00, 0x03, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e,
	0x74, 0x78, 0x74, 0x00, 0x03, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
}

func TestNewCompressReader(t *testing.T) {
	tests := []struct {
		name    string
		r       io.ReadCloser
		wantNil bool
		wantErr bool
	}{
		{
			name:    "Test OK",
			r:       io.NopCloser(bytes.NewReader(gzipData)),
			wantNil: false,
			wantErr: false,
		},
		{
			name:    "Test with error",
			r:       &closableReadableMock{err: fmt.Errorf("some error")},
			wantNil: true,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCompressReader(tt.r)
			require.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantNil, got == nil)
		})
	}
}

func Test_compressReader_Close(t *testing.T) {
	tests := []struct {
		name    string
		r       io.ReadCloser
		zr      io.ReadCloser
		wantErr bool
	}{
		{
			name:    "Test OK",
			r:       &closableReadableMock{},
			zr:      &closableReadableMock{},
			wantErr: false,
		},
		{
			name:    "Test error with reader",
			r:       &closableReadableMock{err: fmt.Errorf("some error")},
			zr:      &closableReadableMock{},
			wantErr: true,
		},
		{
			name:    "Test error with z reader",
			r:       &closableReadableMock{},
			zr:      &closableReadableMock{err: fmt.Errorf("some error")},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &compressReader{r: tt.r, zr: tt.zr}
			assert.Equal(t, tt.wantErr, c.Close() != nil)
		})
	}
}

func Test_compressReader_Read(t *testing.T) {
	tests := []struct {
		name    string
		zr      io.ReadCloser
		p       []byte
		wantN   int
		wantErr bool
	}{
		{
			name:    "Test OK",
			zr:      &closableReadableMock{},
			p:       []byte("test"),
			wantN:   len([]byte("test")),
			wantErr: false,
		},
		{
			name:    "Test with err on read",
			zr:      &closableReadableMock{err: fmt.Errorf("some error")},
			p:       []byte("test"),
			wantN:   len([]byte("test")),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &compressReader{zr: tt.zr}
			gotN, err := c.Read(tt.p)
			require.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantN, gotN)
		})
	}
}
