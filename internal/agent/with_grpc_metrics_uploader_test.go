package agent

import (
	"crypto/tls"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
)

type tlsProvider struct {
	c   *tls.Config
	err error
}

func (p *tlsProvider) Fetch() (*tls.Config, error) {
	return p.c, p.err
}

func TestWithGRPCMetricsUploader(t *testing.T) {
	tests := []struct {
		name        string
		tlsProvider tlsConfigProvider
		host        string
		wantCount   int
	}{
		{
			name:      "Test with empty host",
			wantCount: 0,
		},
		{
			name:        "Test with error on fetch tls config",
			tlsProvider: &tlsProvider{c: nil, err: fmt.Errorf("test")},
			host:        "test",
			wantCount:   1,
		},
		{
			name:        "Test with tls config",
			tlsProvider: &tlsProvider{c: &tls.Config{}},
			host:        "test",
			wantCount:   1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPCMetricsUploader(nil, tt.tlsProvider)
			a := agent{c: &config.Config{Common: config.CommonConfig{GRPCHost: tt.host}}}
			opt(&a)
			assert.Equal(t, tt.wantCount, len(a.processors))
		})
	}
}
