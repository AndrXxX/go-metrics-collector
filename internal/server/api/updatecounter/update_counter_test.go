package updatecounter

import (
	"context"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/AndrXxX/go-metrics-collector/internal/server/repositories/memory"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/counterupdater"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateCounterHandler(t *testing.T) {
	type want struct {
		statusCode int
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
			request: "/update/counter/test/10",
			vars:    map[string]string{vars.Metric: "test", vars.Value: "10"},
			method:  http.MethodPost,
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name:    "StatusBadRequest test",
			request: "/update/counter/test/aaa",
			vars:    map[string]string{vars.Metric: "test", vars.Value: "aaa"},
			method:  http.MethodPost,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}
	storage := memory.New[int64]()
	updater := New(counterupdater.New(&storage))
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.request, nil)
			ctx := chi.NewRouteContext()
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))
			for k, v := range test.vars {
				ctx.URLParams.Add(k, v)
			}

			w := httptest.NewRecorder()

			updater.Handle(w, request)
			result := w.Result()

			assert.Equal(t, test.want.statusCode, result.StatusCode)
			if result.Body != nil {
				err := result.Body.Close()
				require.NoError(t, err)
			}
		})
	}
}
