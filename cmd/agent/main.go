package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/AndrXxX/go-metrics-collector/internal/agent"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/configfile"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/configprovider"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/envparser"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/flagsparser"
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
	if err := logger.Initialize(c.Common.LogLevel); err != nil {
		log.Fatal(err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	buildFormatter := buildformatter.BuildFormatter{
		Labels: []string{"Build version", "Build date", "Build commit"},
		Values: []string{buildVersion, buildDate, buildCommit},
	}
	for _, bInfo := range buildFormatter.Format() {
		logger.Log.Info(bInfo)
	}

	a := agent.New(c, agent.WithRuntimeCollector(), agent.WithVmCollector())
	if err := a.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
