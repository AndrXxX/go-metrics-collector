package updatemetrics

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
)

type testUpdater struct {
	m   *models.Metrics
	err error
}

func (u testUpdater) Update(_ context.Context, _ *models.Metrics) (*models.Metrics, error) {
	return u.m, u.err
}

type testFormatter struct {
	str string
	err error
}

func (f testFormatter) Format(_ *models.Metrics) (string, error) {
	return f.str, f.err
}

type tesIdentifier struct {
	m   *models.Metrics
	err error
}

func (i tesIdentifier) Process(_ *http.Request) (*models.Metrics, error) {
	return i.m, i.err
}

func Test_updateMetricsHandler_Handler(t *testing.T) {
	type fields struct {
		u updater
		f formatter
		i identifier
	}
	tests := []struct {
		name     string
		fields   fields
		wantCode int
	}{
		{
			name:     "Test with error on identify",
			fields:   fields{i: tesIdentifier{err: fmt.Errorf("some error")}},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Test with nil metric on identify",
			fields:   fields{i: tesIdentifier{}},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "Test with error on update",
			fields: fields{
				i: tesIdentifier{m: &models.Metrics{}},
				u: testUpdater{err: fmt.Errorf("some error")},
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "Test with error on format",
			fields: fields{
				i: tesIdentifier{m: &models.Metrics{}},
				u: testUpdater{m: &models.Metrics{}},
				f: testFormatter{err: fmt.Errorf("some error")},
			},
			wantCode: http.StatusInternalServerError,
		},
		{
			name: "Test OK",
			fields: fields{
				i: tesIdentifier{m: &models.Metrics{}},
				u: testUpdater{m: &models.Metrics{}},
				f: testFormatter{str: "{}"},
			},
			wantCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := New(tt.fields.u, tt.fields.f, tt.fields.i)
			f := h.Handler()
			r := httptest.NewRequest("", "/test", nil)
			w := httptest.NewRecorder()
			f(w, r)
			assert.Equal(t, tt.wantCode, w.Code)
		})
	}
}
