package main

import (
	fl "flag"
	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
)

func parseFlags(c *config.Config) {
	fl.StringVar(&c.Host, "a", c.Host, "Net address host:port")
	fl.IntVar(&c.StoreInterval, "i", c.StoreInterval, "Store interval")
	fl.Parse()
}
