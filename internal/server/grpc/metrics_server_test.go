package grpc

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"github.com/AndrXxX/go-metrics-collector/internal/services/utils"

	mp "github.com/AndrXxX/go-metrics-collector/pkg/metricsproto"
)

type testUpdater struct {
	err  error
	list []models.Metrics
}

func (u *testUpdater) UpdateMany(_ context.Context, list []models.Metrics) error {
	u.list = list
	return u.err
}

func TestMetricsServer_UpdateMetrics(t *testing.T) {
	tests := []struct {
		name        string
		u           updater
		in          *mp.UpdateMetricsRequest
		resp        *mp.UpdateMetricsResponse
		wantMetrics []models.Metrics
		wantErr     bool
	}{
		{
			name:        "Test with error on update",
			u:           &testUpdater{err: fmt.Errorf("some error")},
			in:          &mp.UpdateMetricsRequest{},
			resp:        nil,
			wantMetrics: []models.Metrics{},
			wantErr:     true,
		},
		{
			name: "Test with succeed update",
			u:    &testUpdater{},
			in: &mp.UpdateMetricsRequest{
				Metrics: []*mp.Metric{
					{Id: "test1", Type: "type1", Value: 10.1},
					{Id: "test2", Type: "type2", Delta: 10},
				},
			},
			resp: &mp.UpdateMetricsResponse{},
			wantMetrics: []models.Metrics{
				{ID: "test1", MType: "type1", Value: utils.Pointer[float64](10.1)},
				{ID: "test2", MType: "type2", Delta: utils.Pointer[int64](10)},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MetricsServer{Updater: tt.u}
			resp, err := s.UpdateMetrics(context.Background(), tt.in)
			require.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.resp, resp)
		})
	}
}
