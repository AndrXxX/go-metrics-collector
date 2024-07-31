package metricurlbuilder

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/types"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
	"go.uber.org/zap"
	"net/url"
)

type metricURLBuilder struct {
	host string
}

func New(host string) *metricURLBuilder {
	u, err := url.Parse(host)
	if err != nil {
		logger.Log.Error("Error on parse host", zap.String("host", host), zap.Error(err))
		return nil
	}
	if u.Scheme == "" {
		u.Scheme = "http"
	}
	if u.Scheme == "localhost" {
		u.Scheme = "http://localhost"
	}
	return &metricURLBuilder{host: u.String()}
}

func (b *metricURLBuilder) Build(params types.URLParams) string {
	u := fmt.Sprintf("%v/update", b.host)
	if params["endpoint"] != nil {
		u = fmt.Sprintf("%v/%v", b.host, params["endpoint"])
	}
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
