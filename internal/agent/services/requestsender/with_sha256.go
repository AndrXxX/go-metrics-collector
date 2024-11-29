package requestsender

import (
	"bytes"
	"fmt"
	"io"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/requestsender/dto"
)

func WithSHA256(hg hashGenerator, key string) Option {
	return func(p *dto.ParamsDto) error {
		if key == "" {
			return nil
		}
		encoded, err := io.ReadAll(p.Buf)
		if err != nil {
			return fmt.Errorf("error on read encoded data: %w", err)
		}
		p.Buf = bytes.NewBuffer(encoded)
		p.Data = encoded
		p.Headers["HashSHA256"] = hg.Generate(key, encoded)
		return nil
	}
}
