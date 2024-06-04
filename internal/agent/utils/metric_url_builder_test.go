package utils

import "testing"

func TestMetricURLBuilder_BuildURL(t *testing.T) {
	tests := []struct {
		host   string
		params URLParams
		want   string
	}{
		{
			host:   "host",
			params: URLParams{"metricType": "metricType", "metric": "metric", "value": "value"},
			want:   "host/update/metricType/metric/value",
		},
		{
			host:   "host",
			params: URLParams{"metricType": "metricType", "metric": "metric"},
			want:   "host/update/metricType/metric",
		},
		{
			host:   "host",
			params: URLParams{"metricType": "metricType"},
			want:   "host/update/metricType",
		},
		{
			host:   "host",
			params: URLParams{"value": "value"},
			want:   "host/update/value",
		},
		{
			host:   "host",
			params: URLParams{},
			want:   "host/update",
		},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			b := &MetricURLBuilder{host: tt.host}
			if got := b.BuildURL(tt.params); got != tt.want {
				t.Errorf("BuildURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
