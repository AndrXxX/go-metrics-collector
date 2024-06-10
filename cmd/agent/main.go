package main

import (
	"github.com/AndrXxX/go-metrics-collector/internal/agent"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
	"log"
)

func main() {
	c := config.NewConfig()
	parseFlags(c)
	if err := parseEnv(c); err != nil {
		log.Fatal(err)
	}
	if err := agent.Run(c); err != nil {
		log.Fatal(err)
	}
}
