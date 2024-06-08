package main

import (
	"flag"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
)

func parseFlags(c *config.Config) {
	flag.StringVar(&c.Common.Host, "a", c.Common.Host, "Net address host:port")
	flag.Int64Var(&c.Intervals.ReportInterval, "r", c.Intervals.ReportInterval, "Report Interval in seconds")
	flag.Int64Var(&c.Intervals.PollInterval, "p", c.Intervals.PollInterval, "Poll Interval in seconds")
	flag.Parse()
}
