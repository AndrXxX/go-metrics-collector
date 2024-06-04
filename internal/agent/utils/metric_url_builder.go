package utils

import "fmt"

type MetricURLBuilder struct {
	host string
}

func NewMetricURLBuilder(host string) *MetricURLBuilder {
	return &MetricURLBuilder{host: host}
}

func (b *MetricURLBuilder) BuildURL(params URLParams) string {
	url := fmt.Sprintf("%v/update", b.host)
	if params["metricType"] != nil {
		url = fmt.Sprintf("%v/%v", url, params["metricType"])
	}
	if params["metric"] != nil {
		url = fmt.Sprintf("%v/%v", url, params["metric"])
	}
	if params["value"] != nil {
		url = fmt.Sprintf("%v/%v", url, params["value"])
	}
	return url
}
