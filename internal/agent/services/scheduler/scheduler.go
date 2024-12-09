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
func (s *intervalScheduler) AddProcessor(p Processor, interval time.Duration) {
	s.processors = append(s.processors, processorItem{p: p, item: item{interval: interval}})
}

// AddCollector добавляет обработчик для сбора метрик
func (s *intervalScheduler) AddCollector(c Collector, interval time.Duration) {
	s.collectors = append(s.collectors, collectorItem{c: c, item: item{interval: interval}})
}

// Run запускает планировщик
func (s *intervalScheduler) Run() error {
	if s.running.Load() {
		return errors.New("already running")
	}
	logger.Log.Info("Scheduler running")
	s.running.Store(true)
	for {
		chs := s.collect()
		if len(chs) > 0 {
			s.process(s.fanIn(chs))
		}

		if s.stopping.Load() {
			s.stopping.Store(false)
			s.running.Store(false)
			return nil
		} else if len(s.collectors) > 0 || len(s.processors) > 0 {
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
		s.stopping.Store(true)
		for {
			if !s.running.Load() {
				break
			}
			time.Sleep(s.sleepInterval)
		}
		logger.Log.Info("Scheduler stopped")
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

func (s *intervalScheduler) process(ch <-chan dto.MetricsDto) {
	for i := range s.processors {
		if !s.processors[i].canExecute() {
			continue
		}
		s.processors[i].start()
		s.wg.Add(1)
		go func(id int, ch <-chan dto.MetricsDto) {
			err := s.processors[id].p.Process(ch)
			s.processors[id].finish()
			if err != nil {
				logger.Log.Error("Error on process", zap.Error(err))
			}
			s.wg.Done()
		}(i, ch)
	}
}

func (s *intervalScheduler) collect() []chan dto.MetricsDto {
	channels := make([]chan dto.MetricsDto, 0)
	for i := range s.collectors {
		if !s.collectors[i].canExecute() {
			continue
		}
		s.collectors[i].start()
		s.wg.Add(1)
		ch := make(chan dto.MetricsDto)
		channels = append(channels, ch)
		go func(id int) {
			err := s.collectors[id].c.Collect(ch)
			s.collectors[id].finish()
			if err != nil {
				logger.Log.Error("Error on collect", zap.Error(err))
			}
			s.wg.Done()
		}(i)
	}
	return channels
}

func (s *intervalScheduler) fanIn(chs []chan dto.MetricsDto) chan dto.MetricsDto {
	finalCh := make(chan dto.MetricsDto)
	var wg sync.WaitGroup
	for _, ch := range chs {
		wg.Add(1)
		go func(curChan chan dto.MetricsDto) {
			defer wg.Done()
			for data := range curChan {
				finalCh <- data
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(finalCh)
	}()
	return finalCh
}

// NewIntervalScheduler возвращает планировщик, управляющий сборщиками и обработчиками
func NewIntervalScheduler(sleepInterval time.Duration) *intervalScheduler {
	return &intervalScheduler{
		collectors:    []collectorItem{},
		processors:    []processorItem{},
		sleepInterval: sleepInterval,
	}
}
