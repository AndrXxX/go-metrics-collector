package metricsuploader

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AndrXxX/go-metrics-collector/internal/services/utils"
	mp "github.com/AndrXxX/go-metrics-collector/pkg/metricsproto"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/dto"
)

type grpcTestUpdater struct {
	err error
}

func (u *grpcTestUpdater) Update(_ context.Context, _ []*mp.Metric) error {
	return u.err
}

func Test_grpcMetricsUploader_Process(t *testing.T) {
	tests := []struct {
		name    string
		u       grpcMetricsUpdater
		result  dto.MetricsDto
		wantErr bool
	}{
		{
			name:    "Test with empty list",
			result:  dto.MetricsDto{},
			wantErr: false,
		},
		{
			name: "Test with error on update",
			u:    &grpcTestUpdater{fmt.Errorf("some error")},
			result: func() dto.MetricsDto {
				v := dto.NewMetricsDto()
				v.Set(dto.JSONMetrics{})
				return *v
			}(),
			wantErr: true,
		},
		{
			name: "Test with succeed update",
			u:    &grpcTestUpdater{},
			result: func() dto.MetricsDto {
				v := dto.NewMetricsDto()
				v.Set(dto.JSONMetrics{})
				return *v
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewGRPCUploader(tt.u)

			results := make(chan dto.MetricsDto, 1)
			wg := &sync.WaitGroup{}
			wg.Add(1)
			go func() {
				err := c.Process(results)
				require.Equal(t, tt.wantErr, err != nil)
				wg.Done()
			}()
			wg.Add(1)
			go func() {
				results <- tt.result
				time.Sleep(100 * time.Millisecond)
				close(results)
				wg.Done()
			}()
			wg.Wait()
		})
	}
}

func Test_grpcMetricsUploader_convert(t *testing.T) {
	tests := []struct {
		name   string
		result dto.MetricsDto
		want   []*mp.Metric
	}{
		{
			name:   "Test with empty list",
			result: dto.MetricsDto{},
			want:   []*mp.Metric{},
		},
		{
			name: "Test with two examples",
			result: func() dto.MetricsDto {
				v := dto.NewMetricsDto()
				v.Set(dto.JSONMetrics{ID: "test1", MType: "type1", Delta: utils.Pointer[int64](55)})
				v.Set(dto.JSONMetrics{ID: "test2", MType: "type2", Value: utils.Pointer[float64](5.1)})
				return *v
			}(),
			want: []*mp.Metric{
				{Id: "test1", Type: "type1", Delta: 55},
				{Id: "test2", Type: "type2", Value: 5.1},
			},
		},
	}
	sortFunc := func(e, e2 *mp.Metric) int {
		return strings.Compare(strings.ToLower(e.Id), strings.ToLower(e2.Id))
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewGRPCUploader(nil)
			res := c.convert(tt.result)
			slices.SortFunc(res, sortFunc)
			slices.SortFunc(tt.want, sortFunc)
			assert.EqualExportedValues(t, tt.want, res)
		})
	}
}
