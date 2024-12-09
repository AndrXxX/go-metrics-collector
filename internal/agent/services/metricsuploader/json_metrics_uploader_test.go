package metricsuploader

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/dto"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/metricurlbuilder"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/requestsender"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/contenttypes"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/services/utils"
)

func jsonEncode(v any) []byte {
	encoded, _ := json.Marshal(v)
	return encoded
}

func TestNewJSONUploader(t *testing.T) {
	tests := []struct {
		name string
		rs   requestSender
		ub   urlBuilder
		want *jsonMetricsUploader
	}{
		{
			name: "Test NewJSONUploader #1 (Alloc)",
			rs:   requestsender.New(http.DefaultClient),
			ub:   metricurlbuilder.New(""),
			want: &jsonMetricsUploader{
				rs:              requestsender.New(http.DefaultClient),
				ub:              metricurlbuilder.New(""),
				repeatIntervals: []time.Duration{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewJSONUploader(tt.rs, tt.ub, []time.Duration{}))
		})
	}
}

func Test_jsonMetricsUploader_execute(t *testing.T) {
	tests := []struct {
		name      string
		list      []dto.JSONMetrics
		checkFunc func(url string, ct string, data []byte) error
		data      []byte
		wantErr   bool
	}{
		{
			name:    "Test with empty list",
			list:    []dto.JSONMetrics{},
			wantErr: false,
		},
		{
			name: "Test with Gauge Alloc: 0.1",
			list: []dto.JSONMetrics{
				{ID: metrics.Alloc, MType: metrics.Gauge, Value: utils.Pointer[float64](0.1)},
			},
			wantErr: false,
		},
		{
			name: "Test with err",
			list: []dto.JSONMetrics{
				{ID: metrics.Alloc},
			},
			checkFunc: func(_ string, _ string, _ []byte) error {
				return fmt.Errorf("test error")
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := testSender{checkFunc: tt.checkFunc}
			c := NewJSONUploader(rs, metricurlbuilder.New("host"), []time.Duration{})
			result := dto.NewMetricsDto()
			for _, v := range tt.list {
				result.Set(v)
			}
			assert.Equal(t, tt.wantErr, c.execute(*result) != nil)
		})
	}
}

func Test_jsonMetricsUploader_Process(t *testing.T) {
	tests := []struct {
		name      string
		checkFunc func(url string, ct string, data []byte) error
		wantErr   bool
	}{
		{
			name:    "Test OK",
			wantErr: false,
		},
		{
			name: "Test with error",
			checkFunc: func(url string, _ string, _ []byte) error {
				return fmt.Errorf("test error")
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &jsonMetricsUploader{
				rs: testSender{checkFunc: tt.checkFunc},
				ub: metricurlbuilder.New("host"),
			}
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
				v := *dto.NewMetricsDto()
				v.Set(dto.JSONMetrics{ID: "test", MType: metrics.Gauge, Value: func() *float64 {
					val := 0.1
					return &val
				}()})
				results <- v
				time.Sleep(100 * time.Millisecond)
				close(results)
				wg.Done()
			}()
			wg.Wait()
		})
	}
}

func Test_jsonMetricsUploader_send(t *testing.T) {
	tests := []struct {
		name      string
		metric    dto.JSONMetrics
		checkFunc func(url string, ct string, data []byte) error
		data      []byte
		wantErr   bool
	}{
		{
			name:   "Test with Gauge Alloc: 0.1",
			metric: dto.JSONMetrics{ID: metrics.Alloc, MType: metrics.Gauge, Value: utils.Pointer[float64](0.1)},
			checkFunc: func(url string, ct string, data []byte) error {
				assert.Equal(t, "http://host/update", url)
				assert.Equal(t, contenttypes.ApplicationJSON, ct)
				expectedData := jsonEncode(dto.JSONMetrics{ID: metrics.Alloc, MType: metrics.Gauge, Value: utils.Pointer[float64](0.1)})
				assert.Equal(t, expectedData, data)
				return nil
			},
			wantErr: false,
		},
		{
			name:   "Test with err",
			metric: dto.JSONMetrics{ID: metrics.Alloc},
			checkFunc: func(_ string, _ string, _ []byte) error {
				return fmt.Errorf("test error")
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := testSender{checkFunc: tt.checkFunc}
			c := NewJSONUploader(rs, metricurlbuilder.New("host"), []time.Duration{50 * time.Millisecond})
			assert.Equal(t, tt.wantErr, c.send(tt.metric) != nil)
		})
	}
}

func Test_jsonMetricsUploader_sendMany(t *testing.T) {
	tests := []struct {
		name      string
		list      []dto.JSONMetrics
		checkFunc func(url string, ct string, data []byte) error
		data      []byte
		wantErr   bool
	}{
		{
			name:    "Test with empty list",
			list:    []dto.JSONMetrics{},
			wantErr: false,
		},
		{
			name: "Test with Gauge Alloc: 0.1",
			list: []dto.JSONMetrics{
				{ID: metrics.Alloc, MType: metrics.Gauge, Value: utils.Pointer[float64](0.1)},
			},
			checkFunc: func(url string, ct string, data []byte) error {
				assert.Equal(t, "http://host/updates", url)
				assert.Equal(t, contenttypes.ApplicationJSON, ct)
				expectedData := jsonEncode([]dto.JSONMetrics{
					{ID: metrics.Alloc, MType: metrics.Gauge, Value: utils.Pointer[float64](0.1)},
				})
				assert.Equal(t, expectedData, data)
				return nil
			},
			wantErr: false,
		},
		{
			name: "Test with err",
			list: []dto.JSONMetrics{
				{ID: metrics.Alloc},
			},
			checkFunc: func(_ string, _ string, _ []byte) error {
				return fmt.Errorf("test error")
			},
			wantErr: true,
		},
		{
			name: "Test with err on first send and ok on second",
			list: []dto.JSONMetrics{
				{ID: metrics.Alloc},
			},
			checkFunc: func() func(_ string, _ string, _ []byte) error {
				first := true
				return func(_ string, _ string, _ []byte) error {
					if first {
						first = false
						return fmt.Errorf("test error")
					}
					return nil
				}
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := testSender{checkFunc: tt.checkFunc}
			c := NewJSONUploader(rs, metricurlbuilder.New("host"), []time.Duration{50 * time.Millisecond})
			assert.Equal(t, tt.wantErr, c.sendMany(tt.list) != nil)
		})
	}
}
