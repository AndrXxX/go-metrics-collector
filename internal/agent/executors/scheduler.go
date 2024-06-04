package executors

import (
	"errors"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/metrics"
	"time"
)

const sleepInterval = 1

type item struct {
	e            Executors
	interval     time.Duration
	lastExecuted time.Time
}

type intervalScheduler struct {
	list    []item
	running bool
}

func (s *intervalScheduler) Add(e Executors, interval time.Duration) {
	s.list = append(s.list, item{e: e, interval: interval})
}

func (s *intervalScheduler) Run(m metrics.Metrics) error {
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
		time.Sleep(sleepInterval * time.Second)
	}
}

func canExecute(i item) bool {
	return time.Since(i.lastExecuted) >= i.interval
}

func NewIntervalScheduler() *intervalScheduler {
	return &intervalScheduler{
		list: make([]item, 0),
	}
}
