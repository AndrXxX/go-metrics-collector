package counter

import (
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/AndrXxX/go-metrics-collector/internal/repositories/memstorage"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
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
			name:    "StatusMethodNotAllowed test",
			request: "/update/counter/test/10",
			vars:    map[string]string{vars.METRIC: "test", vars.VALUE: "10.1"},
			method:  http.MethodGet,
			want: want{
				statusCode:  http.StatusMethodNotAllowed,
				contentType: "text/plain",
			},
		},
		{
			name:    "StatusNotFound test",
			request: "/update/counter/",
			vars:    map[string]string{},
			method:  http.MethodGet,
			want: want{
				statusCode:  http.StatusMethodNotAllowed,
				contentType: "text/plain",
			},
		},
		{
			name:    "StatusOK test",
			request: "/update/counter/test/10",
			vars:    map[string]string{vars.METRIC: "test", vars.VALUE: "10"},
			method:  http.MethodPost,
			want: want{
				statusCode:  http.StatusOK,
				contentType: "text/plain",
			},
		},
		{
			name:    "StatusBadRequest test",
			request: "/update/counter/test/dsff",
			vars:    map[string]string{vars.METRIC: "test", vars.VALUE: "dsff"},
			method:  http.MethodPost,
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain",
			},
		},
	}
	storage := memstorage.New()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.request, nil)
			request = mux.SetURLVars(request, test.vars)

			w := httptest.NewRecorder()
			Handler(&storage)(w, request)
			result := w.Result()

			assert.Equal(t, test.want.statusCode, result.StatusCode)
			assert.Equal(t, test.want.contentType, result.Header.Get("Content-Type"))
			err := result.Body.Close()
			require.NoError(t, err)
		})
	}
}
