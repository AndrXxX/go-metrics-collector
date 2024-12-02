package metricsuploader

import (
	"errors"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/dto"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/metricurlbuilder"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/requestsender"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
)

type testSender struct {
	checkFunc func(url string, contentType string, data []byte) error
}

func (s testSender) Post(url string, contentType string, data []byte) error {
	if s.checkFunc != nil {
		return s.checkFunc(url, contentType, data)
	}
	return nil
}

func TestNewPlainTextUploader(t *testing.T) {
	tests := []struct {
		name string
		rs   requestSender
		ub   urlBuilder
		want *plainTextMetricsUploader
	}{
		{
			name: "Test NewPlainTextUploader plainTextMetricsUploader #1 (Alloc)",
			rs:   requestsender.New(http.DefaultClient),
			ub:   metricurlbuilder.New(""),
			want: &plainTextMetricsUploader{
				rs: requestsender.New(http.DefaultClient),
				ub: metricurlbuilder.New(""),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewPlainTextUploader(tt.rs, tt.ub), "NewPlainTextUploader(%v)", tt.rs)
		})
	}
}

func Test_plainTextMetricsUploader_execute(t *testing.T) {
	type metricsStr struct {
		Gauge   map[string]float64
		Counter map[string]int64
	}
	tests := []struct {
		name   string
		result metricsStr
		url    string
	}{
		{
			name: "Test Gauge me.Alloc: 0.1",
			result: metricsStr{
				Gauge: map[string]float64{metrics.Alloc: 0.1},
			},
			url: "http://host/update/gauge/Alloc/0.1",
		},
		{
			name: "Test Gauge me.Alloc: 0.1",
			result: metricsStr{
				Gauge: map[string]float64{metrics.HeapAlloc: 0.6661},
			},
			url: "http://host/update/gauge/HeapAlloc/0.6661",
		},
		{
			name: "Test Counter me.PollCount: 555",
			result: metricsStr{
				Counter: map[string]int64{metrics.PollCount: 555},
			},
			url: "http://host/update/counter/PollCount/555",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := testSender{
				checkFunc: func(url string, _ string, _ []byte) error {
					assert.Equal(t, tt.url, url)
					return nil
				},
			}
			c := NewPlainTextUploader(rs, metricurlbuilder.New("host"))
			result := dto.NewMetricsDto()
			for n, v := range tt.result.Gauge {
				result.Set(dto.JSONMetrics{
					ID:    n,
					MType: metrics.Gauge,
					Value: &v,
				})
			}
			for n, v := range tt.result.Counter {
				result.Set(dto.JSONMetrics{
					ID:    n,
					MType: metrics.Counter,
					Delta: &v,
				})
			}

			assert.NoError(t, c.execute(*result))
		})
	}
}

func Test_plainTextMetricsUploader_Process(t *testing.T) {
	type fields struct {
		rs requestSender
		ub urlBuilder
	}
	tests := []struct {
		name    string
		fields  fields
		results chan dto.MetricsDto
		wantErr bool
	}{
		{
			name:    "Test OK",
			fields:  fields{testSender{}, metricurlbuilder.New("host")},
			results: make(chan dto.MetricsDto, 1),
			wantErr: false,
		},
		{
			name: "Test with error",
			fields: fields{testSender{
				checkFunc: func(url string, _ string, _ []byte) error {
					return errors.New("test error")
				},
			}, metricurlbuilder.New("host")},
			results: make(chan dto.MetricsDto, 1),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &plainTextMetricsUploader{
				rs: tt.fields.rs,
				ub: tt.fields.ub,
			}
			wg := &sync.WaitGroup{}
			wg.Add(1)
			go func() {
				err := c.Process(tt.results)
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
				tt.results <- v
				time.Sleep(100 * time.Millisecond)
				close(tt.results)
				wg.Done()
			}()
			wg.Wait()
		})
	}
}
