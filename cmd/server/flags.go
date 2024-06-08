package main

import (
	"flag"
	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
)

func parseFlags(c *config.Config) {
	flag.StringVar(&c.Host, "a", c.Host, "Net address host:port")
	flag.Parse()
}
