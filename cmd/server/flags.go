package main

import (
	fl "flag"
	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
)

func parseFlags(c *config.Config) {
	fl.StringVar(&c.Host, "a", c.Host, "Net address host:port")
	fl.IntVar(&c.StoreInterval, "i", c.StoreInterval, "Store interval")
	fl.StringVar(&c.FileStoragePath, "f", c.FileStoragePath, "File storage path (full)")
	fl.BoolVar(&c.Restore, "r", c.Restore, "Restore values on start")
	fl.Parse()
}
