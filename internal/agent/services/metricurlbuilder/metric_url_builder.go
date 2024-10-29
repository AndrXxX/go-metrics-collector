package metricurlbuilder

import (
	"fmt"
	"net/url"

	"go.uber.org/zap"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/types"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
)

type metricURLBuilder struct {
	host string
}

// New создает экземпляр сервиса metricURLBuilder для построения ссылок
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

// Build собирает и возвращает ссылку по переданным параметрам
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
