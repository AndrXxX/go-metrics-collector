package logger

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_loggingResponseWriter_Write(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "Test with empty data",
		},
		{
			name: "Test with data `test`",
			data: []byte("test"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &loggingResponseWriter{ResponseWriter: httptest.NewRecorder(), responseData: &responseData{}}
			got, err := w.Write(tt.data)
			require.NoError(t, err)
			assert.Equal(t, len(tt.data), got)
			assert.Equal(t, len(tt.data), w.responseData.size)
		})
	}
}

func Test_loggingResponseWriter_WriteHeader(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{
			name:       "Test with StatusOK",
			statusCode: http.StatusOK,
		},
		{
			name:       "Test with StatusInternalServerError",
			statusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &loggingResponseWriter{ResponseWriter: httptest.NewRecorder(), responseData: &responseData{}}
			w.WriteHeader(tt.statusCode)
			assert.Equal(t, tt.statusCode, w.responseData.status)
		})
	}
}
