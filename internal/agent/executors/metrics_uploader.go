package executors

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/metrics"
	"net/http"
	"time"
)

type metricsUploader struct {
	interval     time.Duration
	lastExecuted time.Time
	config       *config.Config
	result       *metrics.Metrics
}

func (c *metricsUploader) canExecute() bool {
	return time.Since(c.lastExecuted) >= c.interval
}

func (c *metricsUploader) Execute() error {
	if !c.canExecute() {
		return nil
	}
	for metric, value := range c.result.Gauge {
		url := buildUrl(c.config.Common.Host, "gauge", metric, value)
		resp, err := http.Post(url, "text/plain", nil)
		if err != nil {
			return err
		}
		fmt.Println(resp.StatusCode, url)
	}
	for metric, value := range c.result.Counter {
		url := buildUrl(c.config.Common.Host, "gauge", metric, value)
		resp, err := http.Post(url, "text/plain", nil)
		if err != nil {
			return err
		}
		fmt.Println(resp.StatusCode, url)
	}
	c.lastExecuted = time.Now()
	return nil
}

func buildUrl(host string, metricType string, metric string, value any) string {
	return fmt.Sprintf("%v/update/%v/%v/%v", host, metricType, metric, value)
}

func NewUploader(interval time.Duration, config *config.Config, result *metrics.Metrics) Executors {
	return &metricsUploader{
		interval:     interval,
		lastExecuted: time.Now(),
		config:       config,
		result:       result,
	}
}
