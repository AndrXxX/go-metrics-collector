package main

import (
	"github.com/AndrXxX/go-metrics-collector/internal/agent"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
)

func main() {
	c := config.NewConfig()
	parseFlags(c)
	parseEnv(c)
	if err := agent.Run(c); err != nil {
		panic(err)
	}
}
