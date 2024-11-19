package scheduler

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/dto"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
)

type intervalScheduler struct {
	processors    []processorItem
	collectors    []collectorItem
	running       atomic.Bool
	stopping      atomic.Bool
	sleepInterval time.Duration
	wg            sync.WaitGroup
}

// AddProcessor добавляет обработчик для выполнения действий с собранными метриками
func (s *intervalScheduler) AddProcessor(p processor, interval time.Duration) {
	s.processors = append(s.processors, processorItem{p: p, interval: interval})
}

// AddCollector добавляет обработчик для сбора метрик
func (s *intervalScheduler) AddCollector(c collector, interval time.Duration) {
	s.collectors = append(s.collectors, collectorItem{c: c, interval: interval})
}

// Run запускает планировщик
func (s *intervalScheduler) Run() error {
	if s.running.Load() {
		return errors.New("already running")
	}
	logger.Log.Info("Scheduler running")
	s.running.Store(true)
	s.stopping.Store(false)
	for {
		channels := make([]chan dto.MetricsDto, 0)
		for _, c := range s.collectors {
			if !canExecute(c.lastExecuted, c.interval) {
				continue
			}
			s.wg.Add(1)
			ch := make(chan dto.MetricsDto)
			channels = append(channels, ch)
			go func() {
				err := c.c.Collect(ch)
				c.lastExecuted = time.Now()
				if err != nil {
					logger.Log.Error("Error on collect", zap.Error(err))
				}
				s.wg.Done()
			}()
		}
		ch := s.fanIn(channels...)
		for _, p := range s.processors {
			if !canExecute(p.lastExecuted, p.interval) {
				continue
			}
			s.wg.Add(1)
			go func() {
				s.process(ch)
				s.wg.Done()
			}()
		}
		if s.stopping.Load() {
			s.stopping.Store(false)
			s.running.Store(false)
			logger.Log.Info("Scheduler stopped")
			s.wg.Done()
			return nil
		}
		if len(s.collectors) > 0 || len(s.processors) > 0 {
			s.wg.Wait()
		}
		time.Sleep(s.sleepInterval)
	}
}

// Shutdown останавливает планировщик
func (s *intervalScheduler) Shutdown(ctx context.Context) error {
	select {
	default:
		logger.Log.Info("Scheduler shutting down")
		s.wg.Add(1)
		go func() {
			s.stopping.Store(true)
		}()
		s.wg.Wait()
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

func (s *intervalScheduler) process(ch <-chan dto.MetricsDto) {
	for _, p := range s.processors {
		s.wg.Add(1)
		go func() {
			if !canExecute(p.lastExecuted, p.interval) {
				s.wg.Done()
				return
			}
			err := p.p.Process(ch)
			p.lastExecuted = time.Now()
			if err != nil {
				logger.Log.Error("Error on collect", zap.Error(err))
			}
			s.wg.Done()
		}()
	}
}

func (s *intervalScheduler) fanIn(chs ...chan dto.MetricsDto) chan dto.MetricsDto {
	finalCh := make(chan dto.MetricsDto)

	var wg sync.WaitGroup
	for _, ch := range chs {
		chClosure := ch
		wg.Add(1)
		go func() {
			defer wg.Done()
			for data := range chClosure {
				finalCh <- data
			}
		}()
	}

	go func() {
		wg.Wait()
		close(finalCh)
	}()
	return finalCh
}

func canExecute(lastExecuted time.Time, interval time.Duration) bool {
	return time.Since(lastExecuted) >= interval
}

// NewIntervalScheduler возвращает планировщик, управляющий сборщиками и обработчиками
func NewIntervalScheduler(sleepInterval time.Duration) *intervalScheduler {
	return &intervalScheduler{
		collectors:    []collectorItem{},
		processors:    []processorItem{},
		sleepInterval: sleepInterval,
	}
}
