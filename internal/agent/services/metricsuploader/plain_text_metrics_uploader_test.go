package metricsuploader

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/dto"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/metricurlbuilder"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/requestsender"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/mocks"
)

func TestNewUploader(t *testing.T) {
	tests := []struct {
		name string
		rs   *requestsender.RequestSender
		ub   urlBuilder
		want *plainTextMetricsUploader
	}{
		{
			name: "Test NewPlainTextUploader plainTextMetricsUploader #1 (Alloc)",
			rs:   requestsender.New(http.DefaultClient, nil, ""),
			ub:   metricurlbuilder.New(""),
			want: &plainTextMetricsUploader{rs: requestsender.New(http.DefaultClient, nil, ""), ub: metricurlbuilder.New("")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewPlainTextUploader(tt.rs, tt.ub), "NewPlainTextUploader(%v)", tt.rs)
		})
	}
}

func Test_metricsUploader_Execute(t *testing.T) {
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
			rs := requestsender.New(&mocks.MockClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					assert.Equal(t, tt.url, req.URL.String())
					return nil, nil
				},
			}, nil, "")
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

			assert.NoError(t, c.Execute(*result))
		})
	}
}
