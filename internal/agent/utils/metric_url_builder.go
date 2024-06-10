package utils

import (
	"fmt"
	"log"
	"net/url"
)

type MetricURLBuilder struct {
	host string
}

func NewMetricURLBuilder(host string) *MetricURLBuilder {
	u, err := url.Parse(host)
	if err != nil {
		log.Print(err)
		return nil
	}
	if u.Scheme == "" {
		u.Scheme = "http"
	}
	return &MetricURLBuilder{host: u.String()}
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
