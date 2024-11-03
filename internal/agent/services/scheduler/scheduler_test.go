package scheduler

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/dto"
)

type executor struct {
}

func (t *executor) Collect(_ chan<- dto.MetricsDto) error {
	return nil
}

func (t *executor) Process(<-chan dto.MetricsDto) error {
	return nil
}

func TestIntervalScheduler_AddCollector(t *testing.T) {
	type args struct {
		c        collector
		interval time.Duration
	}
	tests := []struct {
		name string
		list []collectorItem
		args args
		want []collectorItem
	}{
		{
			name: "add one item when empty list",
			list: []collectorItem{},
			args: args{c: &executor{}, interval: 1 * time.Second},
			want: []collectorItem{{c: &executor{}, interval: 1 * time.Second}},
		},
		{
			name: "add one item when one item exist",
			list: []collectorItem{{c: &executor{}, interval: 2 * time.Second}},
			args: args{c: &executor{}, interval: 1 * time.Second},
			want: []collectorItem{{c: &executor{}, interval: 2 * time.Second}, {c: &executor{}, interval: 1 * time.Second}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &intervalScheduler{collectors: tt.list}
			s.AddCollector(tt.args.c, tt.args.interval)
			assert.Equal(t, tt.want, s.collectors)
		})
	}
}

func TestIntervalScheduler_AddProcessor(t *testing.T) {
	type args struct {
		p        processor
		interval time.Duration
	}
	tests := []struct {
		name string
		list []processorItem
		args args
		want []processorItem
	}{
		{
			name: "add one item when empty list",
			list: []processorItem{},
			args: args{p: &executor{}, interval: 1 * time.Second},
			want: []processorItem{{p: &executor{}, interval: 1 * time.Second}},
		},
		{
			name: "add one item when one item exist",
			list: []processorItem{{p: &executor{}, interval: 2 * time.Second}},
			args: args{p: &executor{}, interval: 1 * time.Second},
			want: []processorItem{{p: &executor{}, interval: 2 * time.Second}, {p: &executor{}, interval: 1 * time.Second}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &intervalScheduler{processors: tt.list}
			s.AddProcessor(tt.args.p, tt.args.interval)
			assert.Equal(t, tt.want, s.processors)
		})
	}
}

func TestNewIntervalScheduler(t *testing.T) {
	tests := []struct {
		name string
		want *intervalScheduler
	}{
		{
			name: "Simple test",
			want: &intervalScheduler{
				collectors:    []collectorItem{},
				processors:    []processorItem{},
				sleepInterval: 1 * time.Second,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewIntervalScheduler(1*time.Second), "NewIntervalScheduler()")
		})
	}
}

func Test_canExecute(t *testing.T) {
	tests := []struct {
		name string
		i    collectorItem
		want bool
	}{
		{
			name: "interval 2, item executed now",
			i:    collectorItem{c: &executor{}, interval: 2 * time.Second, lastExecuted: time.Now()},
			want: false,
		},
		{
			name: "interval 5, item executed 4 seconds ago",
			i:    collectorItem{c: &executor{}, interval: 5 * time.Second, lastExecuted: time.Now().Add(-4 * time.Second)},
			want: false,
		},
		{
			name: "interval 2, item executed 2 seconds ago",
			i:    collectorItem{c: &executor{}, interval: 2 * time.Second, lastExecuted: time.Now().Add(-2 * time.Second)},
			want: true,
		},
		{
			name: "interval 5, item executed 6 seconds ago",
			i:    collectorItem{c: &executor{}, interval: 5 * time.Second, lastExecuted: time.Now().Add(-6 * time.Second)},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, canExecute(tt.i.lastExecuted, tt.i.interval), "canExecute(%v)", tt.i)
		})
	}
}

func Test_intervalScheduler_Run(t *testing.T) {
	type fields struct {
		processors    []processorItem
		collectors    []collectorItem
		running       bool
		stopping      bool
		sleepInterval time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{},
	}
	for _, tt := range tests {
		var wg sync.WaitGroup
		t.Run(tt.name, func(t *testing.T) {
			// TODO: написать test cases
			is := NewIntervalScheduler(1)
			wg.Add(1)
			go func() {
				assert.Equal(t, tt.wantErr, is.Run() != nil)
				wg.Done()
			}()
			wg.Add(1)
			go func() {
				time.Sleep(100 * time.Millisecond)
				_ = is.Shutdown(context.Background())
				wg.Done()
			}()
			wg.Wait()
		})
	}
}
