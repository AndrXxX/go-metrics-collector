package main

import (
	"context"
	"github.com/AndrXxX/go-metrics-collector/internal/server"
	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/dbprovider"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/envparser"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/storageprovider"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
	"github.com/asaskevich/govalidator"
	"go.uber.org/zap"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	settings := config.NewConfig()
	if err := logger.Initialize(settings.LogLevel); err != nil {
		log.Fatal(err)
	}
	parseFlags(settings)
	if err := envparser.New().Parse(settings); err != nil {
		log.Fatal(err)
	}
	if _, err := govalidator.ValidateStruct(settings); err != nil {
		logger.Log.Fatal(err.Error())
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	db, err := dbprovider.New(settings).DB()
	if err != nil {
		logger.Log.Error("failed to connect to database", zap.Error(err))
	}
	sp := storageprovider.New(settings, db)
	app := server.New(settings, sp.Storage(ctx), db)
	if err := app.Run(ctx); err != nil {
		logger.Log.Fatal(err.Error())
	}
}
