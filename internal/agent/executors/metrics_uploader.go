package executors

import (
	"github.com/AndrXxX/go-metrics-collector/internal/agent/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/utils"
)

type metricsUploader struct {
	rs *utils.RequestSender
}

func (c *metricsUploader) Execute(result metrics.Metrics) error {
	for metric, value := range result.Gauge {
		params := utils.URLParams{"metricType": "gauge", "metric": metric, "value": value}
		_ = c.rs.Post(params, "text/plain")
	}
	for metric, value := range result.Counter {
		params := utils.URLParams{"metricType": "counter", "metric": metric, "value": value}
		_ = c.rs.Post(params, "text/plain")
	}
	return nil
}

func NewUploader(rs *utils.RequestSender) Executors {
	return &metricsUploader{
		rs: rs,
	}
}
