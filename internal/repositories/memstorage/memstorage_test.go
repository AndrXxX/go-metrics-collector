package memstorage

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestMemStorage_SetCounter(t *testing.T) {
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
			s.SetCounter(tt.args.metric, tt.args.value)
			assert.Equal(t, tt.want, tt.counter)
		})
	}
}

func TestMemStorage_GetCounter(t *testing.T) {
	type want struct {
		value int64
		error bool
	}
	tests := []struct {
		name    string
		counter map[string]int64
		metric  string
		want    want
	}{
		{
			name:    "Get gauge `metric` with value 1 if exist 1",
			counter: map[string]int64{"metric": 1},
			metric:  "metric",
			want:    want{value: 1, error: false},
		},
		{
			name:    "Get gauge `metric` with value 5 if not exist",
			counter: map[string]int64{},
			metric:  "metric",
			want:    want{value: 0, error: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{counter: tt.counter}
			value, err := s.GetCounter(tt.metric)
			assert.Equal(t, tt.want.value, value)
			if tt.want.error {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestMemStorage_SetGauge(t *testing.T) {
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
			s.SetGauge(tt.args.metric, tt.args.value)
			assert.Equal(t, tt.want, tt.gauge)
		})
	}
}

func TestMemStorage_GetGauge(t *testing.T) {
	type want struct {
		value float64
		error bool
	}
	tests := []struct {
		name   string
		gauge  map[string]float64
		metric string
		want   want
	}{
		{
			name:   "Get gauge `metric` if exist 1.1",
			gauge:  map[string]float64{"metric": 1.1},
			metric: "metric",
			want:   want{value: 1.1, error: false},
		},
		{
			name:   "Get gauge `metric` if not exist",
			gauge:  map[string]float64{},
			metric: "metric",
			want:   want{value: 0, error: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{gauge: tt.gauge}
			value, err := s.GetGauge(tt.metric)
			assert.Equal(t, tt.want.value, value)
			if tt.want.error {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
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
