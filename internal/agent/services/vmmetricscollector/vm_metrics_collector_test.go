package vmmetricscollector

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/dto"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *collector
	}{
		{
			name: "Test OK",
			want: &collector{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, New())
		})
	}
}

func Test_collector_Collect(t *testing.T) {
	tests := []struct {
		name    string
		include []string
	}{
		{
			name:    "Test TotalMemory",
			include: []string{metrics.TotalMemory},
		},
		{
			name:    "Test FreeMemory",
			include: []string{metrics.FreeMemory},
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := make(chan dto.MetricsDto)
			c := New()
			wg.Add(1)
			go func() {
				err := c.Collect(ch)
				wg.Done()
				assert.NoError(t, err)
			}()
			wg.Add(1)
			go func() {
				for result := range ch {
					for _, v := range tt.include {
						_, ok := result.Get(v)
						assert.True(t, ok)
					}
				}
				wg.Done()
			}()
			wg.Wait()
		})
	}
}

func Test_collector_execute(t *testing.T) {
	tests := []struct {
		name    string
		result  dto.MetricsDto
		include []string
	}{
		{
			name:    "Test TotalMemory, FreeMemory",
			result:  *dto.NewMetricsDto(),
			include: []string{metrics.TotalMemory, metrics.FreeMemory},
		},
		{
			name:   "Test CPUutilization",
			result: *dto.NewMetricsDto(),
			include: []string{
				fmt.Sprintf("%s1", metrics.CPUutilization),
				fmt.Sprintf("%s2", metrics.CPUutilization),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &collector{}
			err := c.execute(tt.result)
			require.NoError(t, err)
			for _, v := range tt.include {
				_, ok := tt.result.Get(v)
				assert.True(t, ok)
			}
		})
	}
}
