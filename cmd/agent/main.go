package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/AndrXxX/go-metrics-collector/internal/agent"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/configprovider"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/envparser"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/flagsparser"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
)

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
	if err := agent.Run(ctx, c); err != nil {
		log.Fatal(err)
	}
}
