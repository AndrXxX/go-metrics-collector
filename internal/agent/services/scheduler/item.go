package scheduler

import (
	"sync/atomic"
	"time"
)

type collectorItem struct {
	item
	c Collector
}

type processorItem struct {
	item
	p Processor
}

type item struct {
	interval     time.Duration
	lastExecuted time.Time
	isExecuting  atomic.Bool
}

func (i *item) canExecute() bool {
	return !i.isExecuting.Load() && time.Since(i.lastExecuted) >= i.interval
}

func (i *item) start() {
	i.isExecuting.Store(true)
}

func (i *item) finish() {
	i.lastExecuted = time.Now()
	i.isExecuting.Store(false)
}
