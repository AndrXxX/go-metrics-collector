package utils

import "fmt"

type MetricURLBuilder struct {
	host string
}

func NewMetricURLBuilder(host string) *MetricURLBuilder {
	return &MetricURLBuilder{host: host}
}

func (b *MetricURLBuilder) BuildURL(params URLParams) string {
	return fmt.Sprintf("%v/update/%v/%v/%v", b.host, params["metricType"], params["metric"], params["value"])
}
