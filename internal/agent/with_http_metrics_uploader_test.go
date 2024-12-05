package agent

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
)

type clProvider struct {
	c   *http.Client
	err error
}

func (p *clProvider) Fetch() (*http.Client, error) {
	return p.c, p.err
}

func TestWithHTTPMetricsUploader(t *testing.T) {
	tests := []struct {
		name      string
		provider  clientProvider
		host      string
		wantCount int
	}{
		{
			name:      "Test with empty host",
			wantCount: 0,
		},
		{
			name:      "Test with error on fetch client",
			provider:  &clProvider{c: nil, err: fmt.Errorf("test")},
			host:      "test",
			wantCount: 0,
		},
		{
			name:      "Test with correct client",
			provider:  &clProvider{c: &http.Client{}},
			host:      "test",
			wantCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithHTTPMetricsUploader(nil, tt.provider)
			a := agent{c: &config.Config{Common: config.CommonConfig{Host: tt.host}}}
			opt(&a)
			assert.Equal(t, tt.wantCount, len(a.processors))
		})
	}
}
