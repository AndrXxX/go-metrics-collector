package scheduler

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/dto"
)

type executor struct {
	err           error
	calledProcess bool
}

func (t *executor) Collect(_ chan<- dto.MetricsDto) error {
	return t.err
}

func (t *executor) Process(ch <-chan dto.MetricsDto) error {
	t.calledProcess = true
	return t.err
}

func Test_intervalScheduler_AddCollector(t *testing.T) {
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

func Test_intervalScheduler_AddProcessor(t *testing.T) {
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
	tests := []struct {
		name       string
		processors []processor
		collectors []collector
		running    bool
		shutdown   bool
		wantErr    bool
	}{
		{
			name:    "Test with error on already run",
			running: true,
			wantErr: true,
		},
		{
			name:       "Test with 1 collector and 1 processor",
			collectors: []collector{&executor{err: fmt.Errorf("test error")}},
			processors: []processor{&executor{}},
			shutdown:   true,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := NewIntervalScheduler(50 * time.Millisecond)
			is.running.Store(tt.running)
			for _, c := range tt.collectors {
				is.AddCollector(c, 75*time.Millisecond)
			}
			for _, p := range tt.processors {
				is.AddProcessor(p, 75*time.Millisecond)
			}
			go func() {
				assert.Equal(t, tt.wantErr, is.Run() != nil)
			}()
			if tt.shutdown {
				go func() {
					time.Sleep(200 * time.Millisecond)
					is.wg.Add(1)
					is.stopping.Store(true)
				}()
			}
			time.Sleep(200 * time.Millisecond)
			is.wg.Wait()
		})
	}
}

func Test_intervalScheduler_Shutdown(t *testing.T) {
	t.Run("Test shutdown by cancel context", func(t *testing.T) {
		s := &intervalScheduler{}
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			assert.Equal(t, true, s.Shutdown(ctx) != nil)
		}()
		cancel()
	})
	t.Run("Test shutdown OK", func(t *testing.T) {
		s := &intervalScheduler{}
		go s.wg.Done()
		assert.Equal(t, false, s.Shutdown(context.Background()) != nil)
	})
}

func Test_intervalScheduler_fanIn(t *testing.T) {
	push := func(ch chan dto.MetricsDto, list ...dto.MetricsDto) {
		for _, v := range list {
			ch <- v
		}
		close(ch)
	}
	tests := []struct {
		name string
		chs  []chan dto.MetricsDto
		want chan dto.MetricsDto
	}{
		{
			name: "Test with zero chan",
			chs:  make([]chan dto.MetricsDto, 0),
			want: make(chan dto.MetricsDto),
		},
		{
			name: "Test with one chan (1 item)",
			chs: func() []chan dto.MetricsDto {
				res := make([]chan dto.MetricsDto, 0)
				ch1 := make(chan dto.MetricsDto)
				res = append(res, ch1)
				go push(ch1, dto.MetricsDto{})
				return res
			}(),
			want: make(chan dto.MetricsDto),
		},
		{
			name: "Test with 2 chan (2 and 1 items)",
			chs: func() []chan dto.MetricsDto {
				res := make([]chan dto.MetricsDto, 0)
				ch1 := make(chan dto.MetricsDto, 2)
				res = append(res, ch1)
				go push(ch1, dto.MetricsDto{}, dto.MetricsDto{})
				ch2 := make(chan dto.MetricsDto)
				res = append(res, ch2)
				go push(ch2, dto.MetricsDto{})
				return res
			}(),
			want: make(chan dto.MetricsDto, 3),
		},
		{
			name: "Test with two chan",
			chs: func() []chan dto.MetricsDto {
				res := make([]chan dto.MetricsDto, 0)
				ch1 := make(chan dto.MetricsDto, 2)
				res = append(res, ch1)
				go push(ch1, dto.MetricsDto{}, dto.MetricsDto{})
				ch2 := make(chan dto.MetricsDto)
				res = append(res, ch2)
				go push(ch2, dto.MetricsDto{})
				return res
			}(),
			want: make(chan dto.MetricsDto, 3),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewIntervalScheduler(0)
			assert.Equal(t, len(tt.want), len(s.fanIn(tt.chs...)))
		})
	}
}

func Test_intervalScheduler_process(t *testing.T) {
	t.Run("Test with one processor", func(t *testing.T) {
		s := NewIntervalScheduler(50 * time.Millisecond)
		excr := executor{err: fmt.Errorf("test error")}
		s.AddProcessor(&excr, 1000*time.Millisecond)
		ch := make(chan dto.MetricsDto, 1)
		s.process(ch)
		s.wg.Wait()
		assert.Equal(t, true, excr.calledProcess)

		time.Sleep(200 * time.Millisecond)
		excr.calledProcess = false
		s.process(ch)
		s.wg.Wait()
		assert.Equal(t, false, excr.calledProcess)
	})
}
