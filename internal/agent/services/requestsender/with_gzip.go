package requestsender

import "github.com/AndrXxX/go-metrics-collector/internal/agent/services/requestsender/dto"

func WithGzip(comp dataCompressor) Option {
	return func(p *dto.ParamsDto) error {
		buf, err := comp.Compress(p.Data)
		if err != nil {
			return err
		}
		p.Buf = buf
		p.Headers["Content-Encoding"] = "gzip"
		p.Headers["Accept-Encoding"] = "gzip"
		return nil
	}
}
