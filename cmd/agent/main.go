package main

import (
	"github.com/AndrXxX/go-metrics-collector/internal/agent"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
	"log"
)

func main() {
	c := config.NewConfig()
	if err := logger.Initialize(c.Common.LogLevel); err != nil {
		log.Fatal(err)
	}
	parseFlags(c)
	if err := parseEnv(c); err != nil {
		log.Fatal(err)
	}
	if err := agent.Run(c); err != nil {
		log.Fatal(err)
	}
}
