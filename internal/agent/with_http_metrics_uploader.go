package agent

import (
	"go.uber.org/zap"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/client"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/compressor"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/metricsuploader"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/metricurlbuilder"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/requestsender"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/tlsconfig"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
)

func WithHTTPMetricsUploader(hg hashGenerator) Option {
	return func(a *agent) {
		if a.c.Common.Host == "" {
			return
		}
		ub := metricurlbuilder.New(a.c.Common.Host)

		httpClient, err := client.Provider{ConfProvider: tlsconfig.Provider{CryptoKeyPath: a.c.Common.CryptoKey}}.Fetch()
		if err != nil {
			logger.Log.Error("failed to fetch client: %w", zap.Error(err))
			return
		}
		rs := requestsender.New(
			httpClient,
			requestsender.WithGzip(compressor.GzipCompressor{}),
			requestsender.WithSHA256(hg, a.c.Common.Key),
			requestsender.WithXRealIP(a.c.Common.Host),
		)
		a.processors.Add(metricsuploader.NewJSONUploader(rs, ub, a.c.Intervals.RepeatIntervals))
	}
}
