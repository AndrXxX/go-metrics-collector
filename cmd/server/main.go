package main

import (
	"github.com/AndrXxX/go-metrics-collector/internal/server"
	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
	"github.com/asaskevich/govalidator"
	"log"
)

func main() {
	settings := config.NewConfig()
	if err := logger.Initialize(settings.LogLevel); err != nil {
		log.Fatal(err)
	}
	parseFlags(settings)
	parseEnv(settings)
	if _, err := govalidator.ValidateStruct(settings); err != nil {
		log.Fatal(err)
	}
	if err := server.Run(settings); err != nil {
		log.Fatal(err)
	}
}
