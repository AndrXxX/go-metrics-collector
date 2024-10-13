package metricurlbuilder

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/types"
)

func TestMetricURLBuilderBuildURL(t *testing.T) {
	tests := []struct {
		host   string
		params types.URLParams
		want   string
	}{
		{
			host:   "localhost:8080",
			params: types.URLParams{"metricType": "metricType", "metric": "metric", "value": "value"},
			want:   "http://localhost:8080/update/metricType/metric/value",
		},
		{
			host:   "http://localhost:8080",
			params: types.URLParams{"metricType": "metricType", "metric": "metric", "value": "value"},
			want:   "http://localhost:8080/update/metricType/metric/value",
		},
		{
			host:   "https://localhost:8080",
			params: types.URLParams{"metricType": "metricType", "metric": "metric", "value": "value"},
			want:   "https://localhost:8080/update/metricType/metric/value",
		},
		{
			host:   "host",
			params: types.URLParams{"metricType": "metricType", "metric": "metric", "value": "value"},
			want:   "http://host/update/metricType/metric/value",
		},
		{
			host:   "host",
			params: types.URLParams{"metricType": "metricType", "metric": "metric"},
			want:   "http://host/update/metricType/metric",
		},
		{
			host:   "host",
			params: types.URLParams{"metricType": "metricType"},
			want:   "http://host/update/metricType",
		},
		{
			host:   "host",
			params: types.URLParams{"value": "value"},
			want:   "http://host/update/value",
		},
		{
			host:   "host",
			params: types.URLParams{},
			want:   "http://host/update",
		},
		{
			host:   "host",
			params: types.URLParams{"endpoint": "endpoint"},
			want:   "http://host/endpoint",
		},
		{
			host:   "http://host",
			params: types.URLParams{},
			want:   "http://host/update",
		},
		{
			host:   "https://host",
			params: types.URLParams{},
			want:   "https://host/update",
		},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			b := New(tt.host)
			assert.Equal(t, tt.want, b.Build(tt.params))
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
		{
			host: ":",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.host, func(t *testing.T) {
			b := New(tt.host)
			assert.Equal(t, tt.want, b)
		})
	}
}
