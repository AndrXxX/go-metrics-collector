package scheduler

import (
	"github.com/AndrXxX/go-metrics-collector/internal/agent/dto"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type testCollector struct {
}

func (t *testCollector) Collect(_ chan<- dto.MetricsDto) error {
	return nil
}

func TestIntervalScheduler_Add(t *testing.T) {
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
			args: args{c: &testCollector{}, interval: 1 * time.Second},
			want: []collectorItem{{c: &testCollector{}, interval: 1 * time.Second}},
		},
		{
			name: "add one item when one item exist",
			list: []collectorItem{{c: &testCollector{}, interval: 2 * time.Second}},
			args: args{c: &testCollector{}, interval: 1 * time.Second},
			want: []collectorItem{{c: &testCollector{}, interval: 2 * time.Second}, {c: &testCollector{}, interval: 1 * time.Second}},
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
				sleepInterval: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewIntervalScheduler(1), "NewIntervalScheduler()")
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
			i:    collectorItem{c: &testCollector{}, interval: 2 * time.Second, lastExecuted: time.Now()},
			want: false,
		},
		{
			name: "interval 5, item executed 4 seconds ago",
			i:    collectorItem{c: &testCollector{}, interval: 5 * time.Second, lastExecuted: time.Now().Add(-4 * time.Second)},
			want: false,
		},
		{
			name: "interval 2, item executed 2 seconds ago",
			i:    collectorItem{c: &testCollector{}, interval: 2 * time.Second, lastExecuted: time.Now().Add(-2 * time.Second)},
			want: true,
		},
		{
			name: "interval 5, item executed 6 seconds ago",
			i:    collectorItem{c: &testCollector{}, interval: 5 * time.Second, lastExecuted: time.Now().Add(-6 * time.Second)},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, canExecute(tt.i.lastExecuted, tt.i.interval), "canExecute(%v)", tt.i)
		})
	}
}
