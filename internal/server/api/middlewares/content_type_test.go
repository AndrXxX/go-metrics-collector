package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AndrXxX/go-metrics-collector/internal/enums/contenttypes"
)

func TestSetContentType(t *testing.T) {
	tests := []struct {
		name string
		ct   string
		want *contentType
	}{
		{
			name: "Test with application/json",
			ct:   contenttypes.ApplicationJSON,
			want: &contentType{ct: contenttypes.ApplicationJSON},
		},
		{
			name: "Test with text/plain",
			ct:   contenttypes.TextPlain,
			want: &contentType{ct: contenttypes.TextPlain},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, SetContentType(tt.ct), "SetContentType(%v)", tt.ct)
		})
	}
}

func Test_contentType_Handler(t *testing.T) {
	tests := []struct {
		name        string
		ct          string
		next        http.Handler
		wantHeaders map[string]string
	}{
		{
			name:        "Test with application/json",
			ct:          contenttypes.ApplicationJSON,
			wantHeaders: map[string]string{"Content-Type": contenttypes.ApplicationJSON},
		},
		{
			name: "Test with text/plain",
			ct:   contenttypes.TextPlain,
			next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("testHeader", "testValue")
			}),
			wantHeaders: map[string]string{"Content-Type": contenttypes.TextPlain, "testHeader": "testValue"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &contentType{ct: tt.ct}
			h := m.Handler(tt.next)
			r := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			for k, v := range tt.wantHeaders {
				assert.Equal(t, v, w.Header().Get(k))
			}

		})
	}
}
