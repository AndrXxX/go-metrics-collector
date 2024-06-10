package executors

import (
	"errors"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/metrics"
	"time"
)

type item struct {
	e            Executors
	interval     time.Duration
	lastExecuted time.Time
}

type IntervalScheduler struct {
	list          []item
	running       bool
	sleepInterval int64
}

func (s *IntervalScheduler) Add(e Executors, interval time.Duration) {
	s.list = append(s.list, item{e: e, interval: interval})
}

func (s *IntervalScheduler) Run(m metrics.Metrics) error {
	if s.running {
		return errors.New("already running")
	}
	s.running = true
	for {
		for _, item := range s.list {
			if !canExecute(item) {
				continue
			}
			err := item.e.Execute(m)
			if err != nil {
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

func NewIntervalScheduler(sleepInterval int64) *IntervalScheduler {
	return &IntervalScheduler{
		list:          []item{},
		sleepInterval: sleepInterval,
	}
}
