package savestoragetask

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testSS struct {
	err     error
	saveCnt int
}

func (s *testSS) Save(_ context.Context) error {
	s.saveCnt++
	return s.err
}

func Test_saveStorageTask_Execute(t1 *testing.T) {
	tests := []struct {
		name        string
		i           time.Duration
		s           *testSS
		ctx         func() (context.Context, context.CancelFunc)
		wantSaveCnt int
	}{
		{
			name: "Test with cancel context",
			i:    50 * time.Millisecond,
			s:    &testSS{},
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 5*time.Millisecond)
			},
			wantSaveCnt: 1,
		},
		{
			name: "Test with two executions",
			i:    50 * time.Millisecond,
			s:    &testSS{err: fmt.Errorf("some error")},
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 90*time.Millisecond)
			},
			wantSaveCnt: 2,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := New(tt.i, tt.s)
			ctx, cancel := tt.ctx()
			defer cancel()
			t.Execute(ctx)
			assert.Equal(t1, tt.wantSaveCnt, tt.s.saveCnt)
		})
	}
}
