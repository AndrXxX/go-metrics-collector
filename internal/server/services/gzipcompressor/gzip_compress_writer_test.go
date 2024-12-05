package gzipcompressor

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type closableWritableMock struct {
	closed bool
	err    error
	b      []byte
}

func (m *closableWritableMock) Close() error {
	m.closed = true
	return m.err
}

func (m *closableWritableMock) Write(b []byte) (n int, err error) {
	m.b = b
	return len(b), m.err
}

func Test_compressWriter_Close(t *testing.T) {
	tests := []struct {
		name    string
		zw      *closableWritableMock
		wantErr bool
	}{
		{
			name:    "Test OK",
			zw:      &closableWritableMock{},
			wantErr: false,
		},
		{
			name:    "Test with error on close",
			zw:      &closableWritableMock{err: fmt.Errorf("test error")},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c := NewCompressWriter(w)
			c.zw = tt.zw
			assert.Equal(t, tt.wantErr, c.Close() != nil)
			assert.True(t, tt.zw.closed)
		})
	}
}

func Test_compressWriter_Header(t *testing.T) {
	tests := []struct {
		name string
		w    http.ResponseWriter
	}{
		{
			name: "Test OK",
			w:    httptest.NewRecorder(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c := NewCompressWriter(w)
			assert.Equal(t, tt.w.Header(), c.Header())
		})
	}
}

func Test_compressWriter_Write(t *testing.T) {
	tests := []struct {
		name    string
		w       http.ResponseWriter
		zw      *closableWritableMock
		p       []byte
		want    int
		wantErr bool
	}{
		{
			name:    "Test OK",
			zw:      &closableWritableMock{},
			p:       []byte("test"),
			want:    len([]byte("test")),
			wantErr: false,
		},
		{
			name:    "Test with error on write",
			zw:      &closableWritableMock{err: fmt.Errorf("test error")},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c := NewCompressWriter(w)
			c.zw = tt.zw

			b, err := c.Write(tt.p)
			require.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, b)
			assert.Equal(t, tt.p, tt.zw.b)
		})
	}
}

func Test_compressWriter_WriteHeader(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantHeader map[string]string
	}{
		{
			name:       "Test StatusOK",
			statusCode: http.StatusOK,
			wantHeader: map[string]string{"Content-Encoding": "gzip"},
		},
		{
			name:       "Test StatusUnauthorized",
			statusCode: http.StatusUnauthorized,
			wantHeader: map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c := NewCompressWriter(w)
			c.WriteHeader(tt.statusCode)

			for k := range tt.wantHeader {
				assert.Equal(t, tt.wantHeader[k], w.Header().Get(k), k)
			}
			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}
