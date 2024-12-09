package filestorage

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
)

type testSS struct {
	err error
}

func (ss *testSS) Save(_ context.Context) error {
	return ss.err
}

func (ss *testSS) Restore(_ context.Context) error {
	return ss.err
}

func TestWithStorageSaver(t *testing.T) {
	tests := []struct {
		name string
		c    *config.Config
		ss   storageSaver
	}{
		{
			name: "Test with error on restore",
			c:    &config.Config{Restore: true},
			ss:   &testSS{err: fmt.Errorf("some error")},
		},
		{
			name: "Test OK",
			c:    &config.Config{},
			ss:   &testSS{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithStorageSaver(tt.c, tt.ss)
			fs := &fileStorage{}
			opt(fs)
			assert.Equal(t, tt.ss, fs.ss)
		})
	}
}
