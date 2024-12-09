package agent

import (
	"bytes"
	"time"

	"go.uber.org/zap"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/compressor"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/metricsuploader"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/metricurlbuilder"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/requestsender"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
)

func WithHTTPMetricsUploader(hg hashGenerator, cp clientProvider) Option {
	return func(a *agent) {
		if a.c.Common.Host == "" {
			return
		}

		httpClient, err := cp.Fetch()
		if err != nil {
			logger.Log.Error("failed to fetch client: %w", zap.Error(err))
			return
		}
		rs := requestsender.New(
			httpClient,
			requestsender.WithGzip(compressor.GzipCompressor{Buff: bytes.NewBuffer(nil)}),
			requestsender.WithSHA256(hg, a.c.Common.Key),
			requestsender.WithXRealIP(a.c.Common.Host),
		)
		ub := metricurlbuilder.New(a.c.Common.Host)
		ri := make([]time.Duration, len(a.c.Intervals.RepeatIntervals))
		for i, interval := range a.c.Intervals.RepeatIntervals {
			ri[i] = time.Duration(interval) * time.Second
		}
		a.processors.Add(metricsuploader.NewJSONUploader(rs, ub, ri))
	}
}
