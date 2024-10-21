package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/asaskevich/govalidator"

	"github.com/AndrXxX/go-metrics-collector/internal/agent"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/envparser"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/flagsparser"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
)

func main() {
	c := config.NewConfig()
	if err := logger.Initialize(c.Common.LogLevel); err != nil {
		log.Fatal(err)
	}
	if err := flagsparser.New().Parse(c); err != nil {
		log.Fatal(err)
	}
	if _, err := govalidator.ValidateStruct(c.Intervals); err != nil {
		log.Fatal(err)
	}
	if _, err := govalidator.ValidateStruct(c.Common); err != nil {
		log.Fatal(err)
	}
	if err := envparser.New().Parse(c); err != nil {
		log.Fatal(err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if err := agent.Run(ctx, c); err != nil {
		log.Fatal(err)
	}
}
