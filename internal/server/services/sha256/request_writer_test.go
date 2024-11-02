package sha256

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AndrXxX/go-metrics-collector/internal/services/hashgenerator"
)

func TestRequestWriter_Header(t *testing.T) {
	tests := []struct {
		name string
		ow   http.ResponseWriter
		want http.Header
	}{
		{
			name: "Test OK",
			ow:   httptest.NewRecorder(),
			want: http.Header{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &RequestWriter{OriginalWriter: tt.ow}
			assert.Equal(t, tt.want, w.Header())
		})
	}
}

func TestRequestWriter_Write(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "Test with data",
			data:    []byte("test data"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ow := httptest.NewRecorder()
			buffer := &bytes.Buffer{}
			w := &RequestWriter{OriginalWriter: ow, Buffer: buffer}
			_, err := w.Write(tt.data)
			require.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.data, ow.Body.Bytes())
			assert.Equal(t, tt.data, buffer.Bytes())
		})
	}
}

func TestRequestWriter_WriteHeader(t *testing.T) {
	hg := hashgenerator.Factory().SHA256()
	tests := []struct {
		name       string
		key        string
		data       []byte
		statusCode int
	}{
		{
			name:       "Test with key test1 and data test_data",
			key:        "test1",
			data:       []byte("test data"),
			statusCode: http.StatusOK,
		},
		{
			name:       "Test with empty key and data test2",
			key:        "",
			data:       []byte("test2"),
			statusCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ow := httptest.NewRecorder()
			buffer := &bytes.Buffer{}
			buffer.Write(tt.data)
			w := &RequestWriter{HG: hg, OriginalWriter: ow, Buffer: buffer, Key: tt.key}
			w.WriteHeader(tt.statusCode)
			assert.Equal(t, hg.Generate(tt.key, tt.data), ow.Header().Get("HashSHA256"))
			assert.Equal(t, tt.statusCode, ow.Code)
		})
	}
}
