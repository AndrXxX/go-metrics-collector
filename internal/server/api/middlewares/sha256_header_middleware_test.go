package middlewares

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AndrXxX/go-metrics-collector/internal/server/services/sha256"
	"github.com/AndrXxX/go-metrics-collector/internal/services/hashgenerator"
)

func Test_sha256HeaderMiddleware_Handler(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		wantHash bool
	}{
		{
			name:     "Test with key test1",
			key:      "test1",
			wantHash: true,
		},
		{
			name:     "Test with empty key",
			wantHash: false,
		},
	}
	hg := hashgenerator.Factory().SHA256()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := AddSHA256HashHeader(hg, tt.key)
			h := m.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				assert.Equal(t, tt.wantHash, w.Header().Get("HashSHA256") != "")
			}))
			h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/test", nil))
		})
	}
}

func Test_sha256HeaderMiddleware_processWriter(t *testing.T) {
	hg := hashgenerator.Factory().SHA256()
	w := httptest.NewRecorder()
	tests := []struct {
		name string
		key  string
		want http.ResponseWriter
	}{
		{
			name: "Test with key test1",
			key:  "test1",
			want: &sha256.RequestWriter{HG: hg, OriginalWriter: w, Buffer: &bytes.Buffer{}, Key: "test1"},
		},
		{
			name: "Test with empty key",
			want: w,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := AddSHA256HashHeader(hg, tt.key)
			assert.Equal(t, tt.want, m.processWriter(w))
		})
	}
}
