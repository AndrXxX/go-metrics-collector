package metricurlbuilder

import (
	"github.com/AndrXxX/go-metrics-collector/internal/agent/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetricURLBuilderBuildURL(t *testing.T) {
	tests := []struct {
		host   string
		params utils.URLParams
		want   string
	}{
		{
			host:   "localhost:8080",
			params: utils.URLParams{"metricType": "metricType", "metric": "metric", "value": "value"},
			want:   "http://localhost:8080/update/metricType/metric/value",
		},
		{
			host:   "http://localhost:8080",
			params: utils.URLParams{"metricType": "metricType", "metric": "metric", "value": "value"},
			want:   "http://localhost:8080/update/metricType/metric/value",
		},
		{
			host:   "https://localhost:8080",
			params: utils.URLParams{"metricType": "metricType", "metric": "metric", "value": "value"},
			want:   "https://localhost:8080/update/metricType/metric/value",
		},
		{
			host:   "host",
			params: utils.URLParams{"metricType": "metricType", "metric": "metric", "value": "value"},
			want:   "http://host/update/metricType/metric/value",
		},
		{
			host:   "host",
			params: utils.URLParams{"metricType": "metricType", "metric": "metric"},
			want:   "http://host/update/metricType/metric",
		},
		{
			host:   "host",
			params: utils.URLParams{"metricType": "metricType"},
			want:   "http://host/update/metricType",
		},
		{
			host:   "host",
			params: utils.URLParams{"value": "value"},
			want:   "http://host/update/value",
		},
		{
			host:   "host",
			params: utils.URLParams{},
			want:   "http://host/update",
		},
		{
			host:   "http://host",
			params: utils.URLParams{},
			want:   "http://host/update",
		},
		{
			host:   "https://host",
			params: utils.URLParams{},
			want:   "https://host/update",
		},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			b := New(tt.host)
			if got := b.Build(tt.params); got != tt.want {
				t.Errorf("Build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMetricURLBuilder(t *testing.T) {
	tests := []struct {
		host string
		want *metricURLBuilder
	}{
		{
			host: "host",
			want: &metricURLBuilder{host: "http://host"},
		},
		{
			host: "http://host",
			want: &metricURLBuilder{host: "http://host"},
		},
		{
			host: "https://host",
			want: &metricURLBuilder{host: "https://host"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.host, func(t *testing.T) {
			b := New(tt.host)
			assert.Equal(t, tt.want, b)
		})
	}
}
