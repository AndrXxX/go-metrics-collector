package scheduler

import (
	"errors"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
	"time"
)

type intervalScheduler struct {
	list          []item
	running       bool
	sleepInterval int64
}

func (s *intervalScheduler) Add(e executor, interval time.Duration) {
	s.list = append(s.list, item{e: e, interval: interval})
}

func (s *intervalScheduler) Run(m metrics.Metrics) error {
	if s.running {
		return errors.New("already running")
	}
	logger.Log.Info("Scheduler running")
	s.running = true
	for {
		for _, item := range s.list {
			if !canExecute(item) {
				continue
			}
			err := item.e.Execute(m)
			if err != nil {
				logger.Log.Info("Scheduler stopped")
				s.running = false
				return err
			}
			item.lastExecuted = time.Now()
		}
		time.Sleep(time.Duration(s.sleepInterval) * time.Second)
	}
}

func canExecute(i item) bool {
	return time.Since(i.lastExecuted) >= i.interval
}

func NewIntervalScheduler(sleepInterval int64) *intervalScheduler {
	return &intervalScheduler{
		list:          []item{},
		sleepInterval: sleepInterval,
	}
}
