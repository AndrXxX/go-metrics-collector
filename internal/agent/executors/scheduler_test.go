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
