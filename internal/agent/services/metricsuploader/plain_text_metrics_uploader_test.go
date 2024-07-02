package metricsuploader

import (
	"github.com/AndrXxX/go-metrics-collector/internal/agent/dto"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/metricurlbuilder"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/requestsender"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
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
			rs:   requestsender.New(http.DefaultClient),
			ub:   metricurlbuilder.New(""),
			want: &plainTextMetricsUploader{rs: requestsender.New(http.DefaultClient), ub: metricurlbuilder.New("")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewPlainTextUploader(tt.rs, tt.ub), "NewPlainTextUploader(%v)", tt.rs)
		})
	}
}

func Test_metricsUploader_Execute(t *testing.T) {
	tests := []struct {
		name   string
		result dto.MetricsDto
		url    string
	}{
		{
			name: "Test Gauge me.Alloc: 0.1",
			result: dto.MetricsDto{
				Gauge: map[string]float64{metrics.Alloc: 0.1},
			},
			url: "http://host/update/gauge/Alloc/0.1",
		},
		{
			name: "Test Gauge me.Alloc: 0.1",
			result: dto.MetricsDto{
				Gauge: map[string]float64{metrics.HeapAlloc: 0.6661},
			},
			url: "http://host/update/gauge/HeapAlloc/0.6661",
		},
		{
			name: "Test Counter me.PollCount: 555",
			result: dto.MetricsDto{
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
			})
			c := NewPlainTextUploader(rs, metricurlbuilder.New("host"))
			assert.NoError(t, c.Execute(tt.result))
		})
	}
}
