package logger

import (
	"go.uber.org/zap"
)

// Log переменная для доступа к логгеру
var Log = zap.NewNop()

// Initialize метод инициализации логгера
func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	Log = zl
	return nil
}
