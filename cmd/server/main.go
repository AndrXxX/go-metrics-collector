package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/AndrXxX/go-metrics-collector/internal/server"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/configprovider"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/dbprovider"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/envparser"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/flagsparser"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/storageprovider"
	"github.com/AndrXxX/go-metrics-collector/internal/services/buildformatter"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
)

var buildVersion string
var buildDate string
var buildCommit string

func main() {
	settings, err := configprovider.New(flagsparser.New(), envparser.New()).Fetch()
	if err != nil {
		log.Fatal(err)
	}
	if iErr := logger.Initialize(settings.LogLevel); iErr != nil {
		log.Fatal(err)
	}

	formatter := buildformatter.BuildFormatter{
		Version: buildVersion,
		Date:    buildDate,
		Commit:  buildCommit,
	}
	if str, fErr := formatter.Format(); fErr != nil {
		logger.Log.Error("Failed to print build info", zap.Error(fErr))
	} else {
		fmt.Println(str)
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
