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
		c.upload("gauge", metric, value)
	}
	for metric, value := range result.Counter {
		c.upload("counter", metric, value)
	}
	c.lastExecuted = time.Now()
	return nil
}

func (c *metricsUploader) upload(metricType string, metric string, value any) {
	url := buildURL(*c.host, metricType, metric, value)
	resp, _ := http.Post(url, "text/plain", nil)
	if resp.Body != nil {
		_ = resp.Body.Close()
	}
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
