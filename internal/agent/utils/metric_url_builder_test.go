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
