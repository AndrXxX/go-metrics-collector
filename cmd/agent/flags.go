package main

import (
	fl "flag"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
)

func parseFlags(c *config.Config) {
	fl.StringVar(&c.Common.Host, "a", c.Common.Host, "Net address host:port")
	fl.StringVar(&c.Common.Key, "k", c.Common.Key, "Hash key")
	fl.Int64Var(&c.Common.RateLimit, "l", c.Common.RateLimit, "Rate Limit")
	fl.Int64Var(&c.Intervals.ReportInterval, "r", c.Intervals.ReportInterval, "Report Interval in seconds")
	fl.Int64Var(&c.Intervals.PollInterval, "p", c.Intervals.PollInterval, "Poll Interval in seconds")
	fl.Parse()
}
