package scheduler

import (
	"errors"
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/dto"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
	"sync"
	"time"
)

type intervalScheduler struct {
	list          []item
	processors    []processorItem
	collectors    []collectorItem
	running       bool
	sleepInterval int64
	cwg           sync.WaitGroup
	pwg           sync.WaitGroup
}

func (s *intervalScheduler) Add(e executor, interval time.Duration) {
	s.list = append(s.list, item{e: e, interval: interval})
}

func (s *intervalScheduler) AddProcessor(p processor, interval time.Duration) {
	s.processors = append(s.processors, processorItem{p: p, interval: interval})
}

func (s *intervalScheduler) AddCollector(c collector, interval time.Duration) {
	s.collectors = append(s.collectors, collectorItem{c: c, interval: interval})
}

func (s *intervalScheduler) Run() error {
	if s.running {
		return errors.New("already running")
	}
	logger.Log.Info("Scheduler running")
	s.running = true
	for {
		results := make(chan dto.MetricsDto)
		for _, c := range s.collectors {
			s.cwg.Add(1)
			go func() {
				if !canExecute(c.lastExecuted, c.interval) {
					s.cwg.Done()
					return
				}
				err := c.c.Collect(results)
				c.lastExecuted = time.Now()
				if err != nil {
					logger.Log.Error(fmt.Sprintf("Error on collect: %s", err.Error()))
				}
				s.cwg.Done()
			}()
		}
		for _, p := range s.processors {
			s.pwg.Add(1)
			go func() {
				if !canExecute(p.lastExecuted, p.interval) {
					s.pwg.Done()
					return
				}
				err := p.p.Process(results)
				p.lastExecuted = time.Now()
				if err != nil {
					logger.Log.Error(fmt.Sprintf("Error on collect: %s", err.Error()))
				}
				s.pwg.Done()
			}()
		}
		s.cwg.Wait()
		close(results)
		s.pwg.Wait()
		time.Sleep(time.Duration(s.sleepInterval) * time.Second)
	}
}

func canExecute(lastExecuted time.Time, interval time.Duration) bool {
	return time.Since(lastExecuted) >= interval
}

func NewIntervalScheduler(sleepInterval int64) *intervalScheduler {
	return &intervalScheduler{
		list:          []item{},
		sleepInterval: sleepInterval,
	}
}
