package metrics

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMetrics(t *testing.T) {
	tests := []struct {
		name string
		want *Metrics
	}{
		{
			name: "Test New MemStorage",
			want: &Metrics{Gauge: map[string]float64{}, Counter: map[string]int64{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMetrics()
			assert.Equal(t, tt.want, m)
		})
	}
}
