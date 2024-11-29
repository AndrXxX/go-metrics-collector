package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasGrantedXRealIPMiddleware(t *testing.T) {
	tests := []struct {
		name          string
		trustedSubnet string
		want          *hasGrantedXRealIPMiddleware
	}{
		{
			name:          "Test with empty trustedSubnet",
			trustedSubnet: "",
			want:          &hasGrantedXRealIPMiddleware{},
		},
		{
			name:          "Test with trustedSubnet 192.168.1.0/24",
			trustedSubnet: "192.168.1.0/24",
			want:          &hasGrantedXRealIPMiddleware{"192.168.1.0/24"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, HasGrantedXRealIPOr403(tt.trustedSubnet))
		})
	}
}

func Test_hasGrantedXRealIPMiddleware_Handler(t *testing.T) {
	tests := []struct {
		name          string
		trustedSubnet string
		ip            string
		next          http.Handler
		wantCode      int
	}{
		{
			name:          "Test with not granted ip",
			trustedSubnet: "192.168.1.0/24",
			ip:            "193.168.1.1",
			wantCode:      http.StatusForbidden,
		},
		{
			name:          "Test with granted ip",
			trustedSubnet: "192.168.1.0/24",
			ip:            "192.168.1.4",
			next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
			wantCode: http.StatusOK,
		},
		{
			name:          "Test with empty trustedSubnet",
			trustedSubnet: "",
			next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
			wantCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &hasGrantedXRealIPMiddleware{tt.trustedSubnet}
			r := httptest.NewRequest(http.MethodGet, "/test", nil)
			r.Header.Set(headerXRealIP, tt.ip)
			w := httptest.NewRecorder()
			h := m.Handler(tt.next)
			h.ServeHTTP(w, r)
			assert.Equal(t, tt.wantCode, w.Code)
		})
	}
}

func Test_hasGrantedXRealIPMiddleware_check(t *testing.T) {
	tests := []struct {
		name          string
		trustedSubnet string
		ip            string
		want          bool
	}{
		{
			name:          "Test with wrong subnet",
			trustedSubnet: "192.168.1.1.0/24",
			want:          true,
		},
		{
			name:          "Test with not granted ip",
			trustedSubnet: "192.168.1.0/24",
			ip:            "193.168.1.1",
			want:          false,
		},
		{
			name:          "Test with granted ip",
			trustedSubnet: "192.168.1.0/24",
			ip:            "192.168.1.4",
			want:          true,
		},
		{
			name:          "Test with empty trustedSubnet",
			trustedSubnet: "",
			want:          true,
		},
		{
			name:          "Test with empty XRealIP",
			trustedSubnet: "192.168.1.0/24",
			ip:            "",
			want:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &hasGrantedXRealIPMiddleware{tt.trustedSubnet}
			r := httptest.NewRequest(http.MethodGet, "/test", nil)
			r.Header.Set(headerXRealIP, tt.ip)
			assert.Equal(t, tt.want, m.check(r))
		})
	}
}
