package main

import (
	"github.com/AndrXxX/go-metrics-collector/internal/server"
	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/logger"
)

func main() {
	settings := config.NewConfig()
	if err := logger.Initialize(settings.LogLevel); err != nil {
		panic(err)
	}
	parseFlags(settings)
	parseEnv(settings)
	if err := server.Run(settings); err != nil {
		panic(err)
	}
}
