package main

import (
	"github.com/AndrXxX/go-metrics-collector/internal/server"
	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
)

func main() {
	settings := config.NewConfig()
	parseFlags(settings)
	if err := server.Run(settings); err != nil {
		panic(err)
	}
}
