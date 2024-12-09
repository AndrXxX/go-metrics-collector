package requestsender

import (
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/requestsender/dto"
)

func WithXRealIP(ip string) Option {
	return func(p *dto.ParamsDto) error {
		p.Headers["X-Real-IP"] = ip
		return nil
	}
}
