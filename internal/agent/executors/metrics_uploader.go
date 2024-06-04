package executors

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/metrics"
	"net/http"
)

type metricsUploader struct {
	host string
}

func (c *metricsUploader) Execute(result metrics.Metrics) error {
	for metric, value := range result.Gauge {
		c.upload("gauge", metric, value)
	}
	for metric, value := range result.Counter {
		c.upload("counter", metric, value)
	}
	return nil
}

func (c *metricsUploader) upload(metricType string, metric string, value any) {
	url := buildURL(c.host, metricType, metric, value)
	resp, _ := http.Post(url, "text/plain", nil)
	if resp.Body != nil {
		_ = resp.Body.Close()
	}
}

func buildURL(host string, metricType string, metric string, value any) string {
	return fmt.Sprintf("%v/update/%v/%v/%v", host, metricType, metric, value)
}

func NewUploader(host string) Executors {
	return &metricsUploader{
		host: host,
	}
}
