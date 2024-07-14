package main

import (
	"github.com/AndrXxX/go-metrics-collector/internal/server"
	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/storageprovider"
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
	sp := storageprovider.New(settings)
	app := server.New(settings, sp.Storage())
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
