package interceptors

import "go.uber.org/zap"

type hashGenerator interface {
	Generate(key string, data []byte) string
}

type logService interface {
	Info(msg string, fields ...zap.Field)
}
