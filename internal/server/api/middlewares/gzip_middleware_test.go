package middlewares

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AndrXxX/go-metrics-collector/internal/server/services/gzipcompressor"
)

func TestCompressGzip(t *testing.T) {
	tests := []struct {
		name string
		want *gzipMiddleware
	}{
		{
			name: "Test OK",
			want: &gzipMiddleware{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, CompressGzip())
		})
	}
}

func Test_gzipMiddleware_Handler(t *testing.T) {
	tests := []struct {
		name         string
		headers      map[string]string
		requestBody  io.Reader
		responseData []byte
		next         func(t *testing.T, tw *httptest.ResponseRecorder) http.Handler
	}{
		{
			name:    "Test for accept encoding gzip",
			headers: map[string]string{"Accept-Encoding": "gzip"},
			next: func(t *testing.T, tw *httptest.ResponseRecorder) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					data := []byte("test data")
					_, wErr := w.Write(data)
					require.NoError(t, wErr)
					require.NoError(t, w.(io.Closer).Close())

					reader, cErr := gzipcompressor.NewCompressReader(tw.Result().Body)
					defer func(Body io.ReadCloser) {
						_ = Body.Close()
					}(tw.Result().Body)
					require.NoError(t, cErr)
					decodedData, rErr := io.ReadAll(reader)
					require.NoError(t, rErr)
					assert.Equal(t, data, decodedData)
				})
			},
		},
		{
			name:    "Test for content encoding gzip",
			headers: map[string]string{"Content-Encoding": "gzip"},
			requestBody: func() io.Reader {
				w := httptest.NewRecorder()
				gw := gzip.NewWriter(w)
				_, _ = gw.Write([]byte("test data"))
				_ = gw.Close()
				encodedData, _ := io.ReadAll(bytes.NewReader(w.Body.Bytes()))
				return bytes.NewReader(encodedData)
			}(),
			next: func(t *testing.T, tw *httptest.ResponseRecorder) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					data, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.Equal(t, []byte("test data"), data)
				})
			},
		},
		{
			name:        "Test for content encoding gzip with error",
			headers:     map[string]string{"Content-Encoding": "gzip"},
			requestBody: nil,
			next: func(t *testing.T, tw *httptest.ResponseRecorder) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, http.StatusInternalServerError, tw.Code)
				})
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &gzipMiddleware{}
			r := httptest.NewRequest(http.MethodGet, "/test", tt.requestBody)
			for k, v := range tt.headers {
				r.Header.Set(k, v)
			}
			w := httptest.NewRecorder()
			h := m.Handler(tt.next(t, w))
			_, _ = w.Write(tt.responseData)
			h.ServeHTTP(w, r)
		})
	}
}
