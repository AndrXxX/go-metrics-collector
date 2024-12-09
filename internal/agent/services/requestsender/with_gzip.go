package requestsender

import (
	"io"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/requestsender/dto"
)

func WithGzip(comp dataCompressor) Option {
	return func(p *dto.ParamsDto) error {
		data, err := io.ReadAll(p.Buf)
		if err != nil {
			return err
		}
		buf, err := comp.Compress(data)
		if err != nil {
			return err
		}
		p.Buf = buf
		p.Headers["Content-Encoding"] = "gzip"
		p.Headers["Accept-Encoding"] = "gzip"
		return nil
	}
}
