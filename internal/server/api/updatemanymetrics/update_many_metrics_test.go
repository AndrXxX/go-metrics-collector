package updatemanymetrics

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"github.com/AndrXxX/go-metrics-collector/internal/services/utils"
)

type testUpdater struct {
	list []models.Metrics
	err  error
}

func (u *testUpdater) UpdateMany(_ context.Context, list []models.Metrics) error {
	u.list = list
	return u.err
}

func Test_updateManyMetricsHandler_Handler(t *testing.T) {
	tests := []struct {
		name     string
		u        *testUpdater
		data     []byte
		wantCode int
		wantList []models.Metrics
	}{
		{
			name:     "Test with error on decode body",
			data:     []byte("func()dfd{{{"),
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Test with error on update",
			u:        &testUpdater{err: fmt.Errorf("some error")},
			data:     []byte("[]"),
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Test OK",
			u:        &testUpdater{},
			data:     []byte("[{\"id\":\"test1\",\"type\":\"type1\",\"value\":10.1}]"),
			wantCode: http.StatusOK,
			wantList: []models.Metrics{{ID: "test1", MType: "type1", Value: utils.Pointer(10.1)}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := New(tt.u)
			f := h.Handler()
			r := httptest.NewRequest("", "/test", bytes.NewReader(tt.data))
			w := httptest.NewRecorder()
			f(w, r)
			assert.Equal(t, tt.wantCode, w.Code)
			if tt.wantList != nil {
				assert.Equal(t, tt.wantList, tt.u.list)
			}
		})
	}
}
