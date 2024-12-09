package agent

import (
	"github.com/AndrXxX/go-metrics-collector/internal/agent/services/vmmetricscollector"
)

func WithVMCollector() Option {
	return func(a *agent) {
		a.collectors.Add(vmmetricscollector.New())
	}
}
