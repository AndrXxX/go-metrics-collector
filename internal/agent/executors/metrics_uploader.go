package executors

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/metrics"
	"net/http"
	"time"
)

type metricsUploader struct {
	interval     time.Duration
	lastExecuted time.Time
	host         *string
}

func (c *metricsUploader) canExecute() bool {
	return time.Since(c.lastExecuted) >= c.interval
}

func (c *metricsUploader) Execute(result metrics.Metrics) error {
	if !c.canExecute() {
		return nil
	}
	for metric, value := range result.Gauge {
		url := buildURL(*c.host, "gauge", metric, value)
		resp, _ := http.Post(url, "text/plain", nil)
		if resp.Body != nil {
			_ = resp.Body.Close()
		}
	}
	for metric, value := range result.Counter {
		url := buildURL(*c.host, "counter", metric, value)
		resp, _ := http.Post(url, "text/plain", nil)
		if resp.Body != nil {
			_ = resp.Body.Close()
		}
	}
	c.lastExecuted = time.Now()
	return nil
}

func buildURL(host string, metricType string, metric string, value any) string {
	return fmt.Sprintf("%v/update/%v/%v/%v", host, metricType, metric, value)
}

func NewUploader(interval time.Duration, host *string) Executors {
	return &metricsUploader{
		interval:     interval,
		lastExecuted: time.Now(),
		host:         host,
	}
}
