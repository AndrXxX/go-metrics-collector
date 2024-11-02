package middlewares

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/AndrXxX/go-metrics-collector/internal/services/hashgenerator"
)

type testReader struct {
	mock.Mock
}

func (e *testReader) Read(p []byte) (n int, err error) {
	args := e.Called(p)
	return args.Int(0), args.Error(1)
}

func TestHasCorrectSHA256HashOr500(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want *hasCorrectSHA256HashMiddleware
	}{
		{
			name: "Test OK",
			key:  "test",
			want: &hasCorrectSHA256HashMiddleware{hashgenerator.Factory().SHA256(), "test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, HasCorrectSHA256HashOr500(hashgenerator.Factory().SHA256(), tt.key))
		})
	}
}

func Test_hasCorrectSHA256HashMiddleware_Handler(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		hash     string
		next     http.Handler
		wantCode int
	}{
		{
			name:     "Test with incorrect hash",
			key:      "test",
			hash:     "test",
			wantCode: http.StatusBadRequest,
		},
		{
			name: "Test with empty hash",
			key:  "test",
			next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
			wantCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &hasCorrectSHA256HashMiddleware{hashgenerator.Factory().SHA256(), tt.key}
			r := httptest.NewRequest(http.MethodGet, "/test", nil)
			r.Header.Set("HashSHA256", tt.hash)
			w := httptest.NewRecorder()
			h := m.Handler(tt.next)
			h.ServeHTTP(w, r)
			assert.Equal(t, tt.wantCode, w.Code)
		})
	}
}

func Test_hasCorrectSHA256HashMiddleware_check(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		hash        string
		requestBody io.Reader
		want        bool
	}{
		{
			name: "Test with empty key",
			want: true,
		},
		{
			name: "Test with empty hash",
			key:  "test",
			want: true,
		},
		{
			name: "Test with error on read body",
			key:  "test",
			hash: "test",
			requestBody: func() io.Reader {
				r := &testReader{}
				r.On("Read", mock.Anything).Return(0, fmt.Errorf("test error"))
				return r
			}(),
			want: false,
		},
		{
			name:        "Test with right key",
			key:         "test",
			hash:        hashgenerator.Factory().SHA256().Generate("test", []byte("test")),
			requestBody: bytes.NewBufferString("test"),
			want:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &hasCorrectSHA256HashMiddleware{hashgenerator.Factory().SHA256(), tt.key}
			r := httptest.NewRequest(http.MethodGet, "/test", tt.requestBody)
			r.Header.Set("HashSHA256", tt.hash)
			assert.Equal(t, tt.want, m.check(r))
		})
	}
}
