package utils

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
	"net/url"
)

type MetricURLBuilder struct {
	host string
}

func NewMetricURLBuilder(host string) *MetricURLBuilder {
	u, err := url.Parse(host)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Error on parse host %s: %s", host, err.Error()))
		return nil
	}
	if u.Scheme == "" {
		u.Scheme = "http"
	}
	if u.Scheme == "localhost" {
		u.Scheme = "http://localhost"
	}
	return &MetricURLBuilder{host: u.String()}
}

func (b *MetricURLBuilder) BuildURL(params URLParams) string {
	u := fmt.Sprintf("%v/update", b.host)
	if params["metricType"] != nil {
		u = fmt.Sprintf("%v/%v", u, params["metricType"])
	}
	if params["metric"] != nil {
		u = fmt.Sprintf("%v/%v", u, params["metric"])
	}
	if params["value"] != nil {
		u = fmt.Sprintf("%v/%v", u, params["value"])
	}
	return u
}
