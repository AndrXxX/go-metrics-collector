package memstorage

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestMemStorage_Counter(t *testing.T) {
	type args struct {
		metric string
		value  int64
	}
	tests := []struct {
		name    string
		counter map[string]int64
		args    args
		want    map[string]int64
	}{
		{
			name:    "Add counter `metric` with value 1 if not exist",
			counter: map[string]int64{},
			args:    args{metric: "metric", value: 1},
			want:    map[string]int64{"metric": 1},
		},
		{
			name:    "Add counter `metric` with value 10 if exist 15",
			counter: map[string]int64{"metric": 15},
			args:    args{metric: "metric", value: 10},
			want:    map[string]int64{"metric": 25},
		},
		{
			name:    "Add counter `metric` with value -1 if exist 15",
			counter: map[string]int64{"metric": 15},
			args:    args{metric: "metric", value: -1},
			want:    map[string]int64{"metric": 14},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{counter: tt.counter}
			s.Counter(tt.args.metric, tt.args.value)
			assert.Equal(t, tt.counter, tt.want)
		})
	}
}

func TestMemStorage_Gauge(t *testing.T) {
	type args struct {
		metric string
		value  float64
	}
	tests := []struct {
		name  string
		gauge map[string]float64
		args  args
		want  map[string]float64
	}{
		{
			name:  "Set gauge `metric` with value 1 if not exist",
			gauge: map[string]float64{},
			args:  args{metric: "metric", value: 1},
			want:  map[string]float64{"metric": 1},
		},
		{
			name:  "Set gauge `metric` with value 10 if exist 15",
			gauge: map[string]float64{"metric": 15},
			args:  args{metric: "metric", value: 10},
			want:  map[string]float64{"metric": 10},
		},
		{
			name:  "Set gauge `metric` with value -1 if exist 15",
			gauge: map[string]float64{"metric": 15},
			args:  args{metric: "metric", value: -1},
			want:  map[string]float64{"metric": -1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{gauge: tt.gauge}
			s.Gauge(tt.args.metric, tt.args.value)
			assert.Equal(t, tt.gauge, tt.want)
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want MemStorage
	}{
		{
			name: "Test New MemStorage",
			want: MemStorage{gauge: make(map[string]float64), counter: make(map[string]int64)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
