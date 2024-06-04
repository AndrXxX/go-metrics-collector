package executors

import (
	"github.com/AndrXxX/go-metrics-collector/internal/agent/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/utils"
	me "github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/mocks"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestNewUploader(t *testing.T) {
	tests := []struct {
		name string
		rs   *utils.RequestSender
		want Executors
	}{
		{
			name: "Test New metricsUploader #1 (Alloc)",
			rs:   utils.NewRequestSender(utils.NewMetricURLBuilder(""), http.DefaultClient),
			want: &metricsUploader{rs: utils.NewRequestSender(utils.NewMetricURLBuilder(""), http.DefaultClient)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewUploader(tt.rs), "NewUploader(%v)", tt.rs)
		})
	}
}

func Test_metricsUploader_Execute(t *testing.T) {
	tests := []struct {
		name   string
		result metrics.Metrics
		url    string
	}{
		{
			name: "Test Gauge me.Alloc: 0.1",
			result: metrics.Metrics{
				Gauge: map[string]float64{me.Alloc: 0.1},
			},
			url: "host/update/gauge/Alloc/0.1",
		},
		{
			name: "Test Gauge me.Alloc: 0.1",
			result: metrics.Metrics{
				Gauge: map[string]float64{me.HeapAlloc: 0.6661},
			},
			url: "host/update/gauge/HeapAlloc/0.6661",
		},
		{
			name: "Test Counter me.PollCount: 555",
			result: metrics.Metrics{
				Counter: map[string]int64{me.PollCount: 555},
			},
			url: "host/update/counter/PollCount/555",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := utils.NewRequestSender(utils.NewMetricURLBuilder("host"), &mocks.MockClient{
				PostFunc: func(url, contentType string, body io.Reader) (*http.Response, error) {
					assert.Equal(t, tt.url, url)
					return nil, nil
				},
			})
			c := &metricsUploader{rs: rs}
			assert.NoError(t, c.Execute(tt.result))
		})
	}
}
