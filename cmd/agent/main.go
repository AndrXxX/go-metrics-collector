package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/AndrXxX/go-metrics-collector/internal/agent"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/configprovider"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/envparser"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/flagsparser"
	"github.com/AndrXxX/go-metrics-collector/internal/services/buildformatter"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
)

var buildVersion string
var buildDate string
var buildCommit string

func main() {
	c, err := configprovider.New(flagsparser.New(), envparser.New()).Fetch()
	if err != nil {
		log.Fatal(err)
	}
	if err := logger.Initialize(c.Common.LogLevel); err != nil {
		log.Fatal(err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

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

	if err := agent.Run(ctx, c); err != nil {
		log.Fatal(err)
	}
}
