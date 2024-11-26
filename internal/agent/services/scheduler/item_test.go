package scheduler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_canExecute(t *testing.T) {
	tests := []struct {
		name         string
		interval     time.Duration
		lastExecuted time.Time
		want         bool
	}{
		{
			name:         "interval 2, item executed now",
			interval:     2 * time.Second,
			lastExecuted: time.Now(),
			want:         false,
		},
		{
			name:         "interval 5, item executed 4 seconds ago",
			interval:     5 * time.Second,
			lastExecuted: time.Now().Add(-4 * time.Second),
			want:         false,
		},
		{
			name:         "interval 2, item executed 2 seconds ago",
			interval:     2 * time.Second,
			lastExecuted: time.Now().Add(-2 * time.Second),
			want:         true,
		},
		{
			name:         "interval 5, item executed 6 seconds ago",
			interval:     5 * time.Second,
			lastExecuted: time.Now().Add(-6 * time.Second),
			want:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := item{interval: tt.interval, lastExecuted: tt.lastExecuted}
			assert.Equal(t, tt.want, i.canExecute())
		})
	}
}

func Test_item_start_and_finish(t *testing.T) {
	t.Run("Check time", func(t *testing.T) {
		now := time.Now()
		i := &item{}
		i.start()
		i.finish()
		assert.True(t, i.lastExecuted.After(now))
	})
}
