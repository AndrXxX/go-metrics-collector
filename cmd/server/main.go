package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/AndrXxX/go-metrics-collector/internal/server"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/configfile"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/configprovider"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/dbprovider"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/envparser"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/flagsparser"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/storageprovider"
	"github.com/AndrXxX/go-metrics-collector/internal/services/buildformatter"
	"github.com/AndrXxX/go-metrics-collector/internal/services/configpath"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
)

var buildVersion string
var buildDate string
var buildCommit string

func main() {
	cpp := configpath.NewProvider(configpath.WithFlags("c", "config"), configpath.WithEnv())
	c, err := configprovider.New(configfile.Parser{PathProvider: cpp}, flagsparser.New(), envparser.New()).Fetch()
	if err != nil {
		log.Fatal(err)
	}
	if iErr := logger.Initialize(c.LogLevel); iErr != nil {
		log.Fatal(err)
	}

	buildFormatter := buildformatter.BuildFormatter{
		Labels: []string{"Build version", "Build date", "Build commit"},
		Values: []string{buildVersion, buildDate, buildCommit},
	}
	for _, bInfo := range buildFormatter.Format() {
		logger.Log.Info(bInfo)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	db, err := dbprovider.New(c).DB()
	if err != nil {
		logger.Log.Error("failed to connect to database", zap.Error(err))
	}
	sp := storageprovider.New(c, db)
	app := server.New(c, sp.Storage(ctx), db)
	if err := app.Run(ctx); err != nil {
		logger.Log.Fatal(err.Error())
	}
}
