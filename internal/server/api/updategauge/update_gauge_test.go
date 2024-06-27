package updategauge

import (
	"context"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/AndrXxX/go-metrics-collector/internal/server/repositories/memory"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateGaugeHandlerHandle(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name    string
		request string
		vars    map[string]string
		method  string
		want    want
	}{
		{
			name:    "StatusOK test",
			request: "/update/gauge/test/10.1",
			vars:    map[string]string{vars.Metric: "test", vars.Value: "10.1"},
			method:  http.MethodPost,
			want: want{
				statusCode:  http.StatusOK,
				contentType: "text/plain",
			},
		},
		{
			name:    "StatusBadRequest test",
			request: "/update/gauge/test/aaa",
			vars:    map[string]string{vars.Metric: "test", vars.Value: "aaa"},
			method:  http.MethodPost,
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain",
			},
		},
	}
	storage := memory.New[float64]()
	h := New(&storage)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.request, nil)
			rctx := chi.NewRouteContext()
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))
			for k, v := range test.vars {
				rctx.URLParams.Add(k, v)
			}

			w := httptest.NewRecorder()
			h.Handle(w, request)
			result := w.Result()

			assert.Equal(t, test.want.statusCode, result.StatusCode)
			assert.Equal(t, test.want.contentType, result.Header.Get("Content-Type"))
			if result.Body != nil {
				err := result.Body.Close()
				require.NoError(t, err)
			}
		})
	}
}
