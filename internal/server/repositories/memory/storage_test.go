package memory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorageInsertInt64(t *testing.T) {
	type args struct {
		metric string
		value  int64
	}
	tests := []struct {
		name  string
		store map[string]int64
		args  args
		want  map[string]int64
	}{
		{
			name:  "Add counter `metric` with value 1 if not exist",
			store: map[string]int64{},
			args:  args{metric: "metric", value: 1},
			want:  map[string]int64{"metric": 1},
		},
		{
			name:  "Add counter `metric` with value 10 if exist 15",
			store: map[string]int64{"metric": 15},
			args:  args{metric: "metric", value: 10},
			want:  map[string]int64{"metric": 10},
		},
	}
	ctx := context.TODO()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New[int64]()
			s.Insert(ctx, tt.args.metric, tt.args.value)
			assert.Equal(t, tt.want, s.All(ctx))
		})
	}
}

func TestStorageGetInt64(t *testing.T) {
	type want struct {
		value int64
		ok    bool
	}
	tests := []struct {
		name   string
		store  map[string]int64
		metric string
		want   want
	}{
		{
			name:   "Get count `metric` with value 1 if exist 1",
			store:  map[string]int64{"metric": 1},
			metric: "metric",
			want:   want{value: 1, ok: true},
		},
		{
			name:   "Get count `metric` with value 5 if not exist",
			store:  map[string]int64{},
			metric: "metric",
			want:   want{value: 0, ok: false},
		},
	}
	ctx := context.TODO()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New[int64]()
			for k, v := range tt.store {
				s.Insert(ctx, k, v)
			}
			value, ok := s.Get(ctx, tt.metric)
			assert.Equal(t, tt.want.value, value)
			assert.Equal(t, tt.want.ok, ok)
		})
	}
}

func TestStorageAllInt64(t *testing.T) {
	tests := []struct {
		name  string
		store map[string]int64
		want  map[string]int64
	}{
		{
			name:  "Empty store",
			store: map[string]int64{},
			want:  map[string]int64{},
		},
		{
			name:  "Counter `metric` with value 1",
			store: map[string]int64{"metric": 1},
			want:  map[string]int64{"metric": 1},
		},
		{
			name:  "Counter `metric` with value 1 and `metric2` with value 10",
			store: map[string]int64{"metric": 1, "metric2": 10},
			want:  map[string]int64{"metric": 1, "metric2": 10},
		},
	}
	ctx := context.TODO()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New[int64]()
			for k, v := range tt.store {
				s.Insert(ctx, k, v)
			}
			assert.Equal(t, tt.want, s.All(ctx))
		})
	}
}

func TestStorageSetFloat64(t *testing.T) {
	type args struct {
		metric string
		value  float64
	}
	tests := []struct {
		name  string
		store map[string]float64
		args  args
		want  map[string]float64
	}{
		{
			name:  "Set gauge `metric` with value 1 if not exist",
			store: map[string]float64{},
			args:  args{metric: "metric", value: 1},
			want:  map[string]float64{"metric": 1},
		},
		{
			name:  "Set gauge `metric` with value 10 if exist 15",
			store: map[string]float64{"metric": 15},
			args:  args{metric: "metric", value: 10},
			want:  map[string]float64{"metric": 10},
		},
		{
			name:  "Set gauge `metric` with value -1 if exist 15",
			store: map[string]float64{"metric": 15},
			args:  args{metric: "metric", value: -1},
			want:  map[string]float64{"metric": -1},
		},
	}
	ctx := context.TODO()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New[float64]()
			s.Insert(ctx, tt.args.metric, tt.args.value)
			assert.Equal(t, tt.want, s.All(ctx))
		})
	}
}

func TestStorageGetFloat64(t *testing.T) {
	type want struct {
		value float64
		ok    bool
	}
	tests := []struct {
		name   string
		gauge  map[string]float64
		metric string
		want   want
	}{
		{
			name:   "Get store `metric` if exist 1.1",
			gauge:  map[string]float64{"metric": 1.1},
			metric: "metric",
			want:   want{value: 1.1, ok: true},
		},
		{
			name:   "Get store `metric` if not exist",
			gauge:  map[string]float64{},
			metric: "metric",
			want:   want{value: 0, ok: false},
		},
	}
	ctx := context.TODO()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New[float64]()
			for k, v := range tt.gauge {
				s.Insert(ctx, k, v)
			}
			value, ok := s.Get(ctx, tt.metric)
			assert.InDelta(t, tt.want.value, value, 0.0001)
			assert.Equal(t, tt.want.ok, ok)
		})
	}
}

func TestStorageDelete(t *testing.T) {
	tests := []struct {
		name     string
		store    map[string]float64
		toDelete string
		want     map[string]float64
	}{
		{
			name:     "Empty store",
			store:    map[string]float64{},
			toDelete: "test",
			want:     map[string]float64{},
		},
		{
			name:     "Gauge `metric` with value 1.1",
			store:    map[string]float64{"metric": 1.1},
			toDelete: "metric",
			want:     map[string]float64{},
		},
		{
			name:     "Gauge `metric` with value 1.1 and `metric2` with value 10.5",
			store:    map[string]float64{"metric": 1.1, "metric2": 10.5},
			toDelete: "metric2",
			want:     map[string]float64{"metric": 1.1},
		},
	}
	ctx := context.TODO()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New[float64]()
			for k, v := range tt.store {
				s.Insert(ctx, k, v)
			}
			s.Delete(ctx, tt.toDelete)
			assert.Equal(t, tt.want, s.All(ctx))
		})
	}
}

func TestStorageAllFloat64(t *testing.T) {
	tests := []struct {
		name  string
		store map[string]float64
		want  map[string]float64
	}{
		{
			name:  "Empty store",
			store: map[string]float64{},
			want:  map[string]float64{},
		},
		{
			name:  "Gauge `metric` with value 1.1",
			store: map[string]float64{"metric": 1.1},
			want:  map[string]float64{"metric": 1.1},
		},
		{
			name:  "Gauge `metric` with value 1.1 and `metric2` with value 10.5",
			store: map[string]float64{"metric": 1.1, "metric2": 10.5},
			want:  map[string]float64{"metric": 1.1, "metric2": 10.5},
		},
	}
	ctx := context.TODO()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New[float64]()
			for k, v := range tt.store {
				s.Insert(ctx, k, v)
			}
			assert.Equal(t, tt.want, s.All(ctx))
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *storage[any]
	}{
		{
			name: "Test New storage",
			want: &storage[any]{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New[any]()
			assert.Equal(t, tt.want, &s)
		})
	}
}
