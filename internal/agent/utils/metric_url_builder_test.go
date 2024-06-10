package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetricURLBuilder_BuildURL(t *testing.T) {
	tests := []struct {
		host   string
		params URLParams
		want   string
	}{
		{
			host:   "host",
			params: URLParams{"metricType": "metricType", "metric": "metric", "value": "value"},
			want:   "http://host/update/metricType/metric/value",
		},
		{
			host:   "host",
			params: URLParams{"metricType": "metricType", "metric": "metric"},
			want:   "http://host/update/metricType/metric",
		},
		{
			host:   "host",
			params: URLParams{"metricType": "metricType"},
			want:   "http://host/update/metricType",
		},
		{
			host:   "host",
			params: URLParams{"value": "value"},
			want:   "http://host/update/value",
		},
		{
			host:   "host",
			params: URLParams{},
			want:   "http://host/update",
		},
		{
			host:   "http://host",
			params: URLParams{},
			want:   "http://host/update",
		},
		{
			host:   "https://host",
			params: URLParams{},
			want:   "https://host/update",
		},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			b := NewMetricURLBuilder(tt.host)
			if got := b.BuildURL(tt.params); got != tt.want {
				t.Errorf("BuildURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMetricURLBuilder(t *testing.T) {
	tests := []struct {
		host string
		want *MetricURLBuilder
	}{
		{
			host: "host",
			want: &MetricURLBuilder{host: "http://host"},
		},
		{
			host: "http://host",
			want: &MetricURLBuilder{host: "http://host"},
		},
		{
			host: "https://host",
			want: &MetricURLBuilder{host: "https://host"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.host, func(t *testing.T) {
			b := NewMetricURLBuilder(tt.host)
			assert.Equal(t, tt.want, b)
		})
	}
}
