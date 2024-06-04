package executors

import (
	"github.com/AndrXxX/go-metrics-collector/internal/agent/metrics"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type testExecutor struct {
}

func (t *testExecutor) Execute(m metrics.Metrics) error {
	return nil
}

func TestIntervalScheduler_Add(t *testing.T) {
	type args struct {
		e        Executors
		interval time.Duration
	}
	tests := []struct {
		name string
		list []item
		args args
		want []item
	}{
		{
			name: "add one item when empty list",
			list: []item{},
			args: args{e: &testExecutor{}, interval: 1 * time.Second},
			want: []item{{e: &testExecutor{}, interval: 1 * time.Second}},
		},
		{
			name: "add one item when one item exist",
			list: []item{{e: &testExecutor{}, interval: 2 * time.Second}},
			args: args{e: &testExecutor{}, interval: 1 * time.Second},
			want: []item{{e: &testExecutor{}, interval: 2 * time.Second}, {e: &testExecutor{}, interval: 1 * time.Second}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &IntervalScheduler{list: tt.list}
			s.Add(tt.args.e, tt.args.interval)
			assert.Equal(t, tt.want, s.list)
		})
	}
}

func Test_canExecute(t *testing.T) {
	tests := []struct {
		name string
		i    item
		want bool
	}{
		{
			name: "interval 2, item executed now",
			i:    item{e: &testExecutor{}, interval: 2 * time.Second, lastExecuted: time.Now()},
			want: false,
		},
		{
			name: "interval 5, item executed 4 seconds ago",
			i:    item{e: &testExecutor{}, interval: 5 * time.Second, lastExecuted: time.Now().Add(-4 * time.Second)},
			want: false,
		},
		{
			name: "interval 2, item executed 2 seconds ago",
			i:    item{e: &testExecutor{}, interval: 2 * time.Second, lastExecuted: time.Now().Add(-2 * time.Second)},
			want: true,
		},
		{
			name: "interval 5, item executed 6 seconds ago",
			i:    item{e: &testExecutor{}, interval: 5 * time.Second, lastExecuted: time.Now().Add(-6 * time.Second)},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, canExecute(tt.i), "canExecute(%v)", tt.i)
		})
	}
}
