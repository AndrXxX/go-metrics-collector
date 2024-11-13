package envparser

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
)

func Test_parseEnv(t *testing.T) {
	tests := []struct {
		name    string
		config  *config.Config
		env     map[string]string
		want    *config.Config
		wantErr bool
	}{
		{
			name: "Empty env",
			config: &config.Config{
				Common:    config.CommonConfig{Host: "host"},
				Intervals: config.Intervals{PollInterval: 1, ReportInterval: 1},
			},
			env: map[string]string{},
			want: &config.Config{
				Common:    config.CommonConfig{Host: "host"},
				Intervals: config.Intervals{PollInterval: 1, ReportInterval: 1},
			},
			wantErr: false,
		},
		{
			name: "ADDRESS=new-host",
			config: &config.Config{
				Common:    config.CommonConfig{Host: "host"},
				Intervals: config.Intervals{PollInterval: 1, ReportInterval: 1},
			},
			env: map[string]string{"ADDRESS": "new-host"},
			want: &config.Config{
				Common:    config.CommonConfig{Host: "new-host"},
				Intervals: config.Intervals{PollInterval: 1, ReportInterval: 1},
			},
			wantErr: false,
		},
		{
			name: "REPORT_INTERVAL=fff",
			config: &config.Config{
				Common:    config.CommonConfig{Host: "host"},
				Intervals: config.Intervals{PollInterval: 1, ReportInterval: 1},
			},
			env: map[string]string{"REPORT_INTERVAL": "fff"},
			want: &config.Config{
				Common:    config.CommonConfig{Host: "host"},
				Intervals: config.Intervals{PollInterval: 1, ReportInterval: 1},
			},
			wantErr: true,
		},
		{
			name: "REPORT_INTERVAL=2",
			config: &config.Config{
				Common:    config.CommonConfig{Host: "host"},
				Intervals: config.Intervals{PollInterval: 1, ReportInterval: 1},
			},
			env: map[string]string{"REPORT_INTERVAL": "2"},
			want: &config.Config{
				Common:    config.CommonConfig{Host: "host"},
				Intervals: config.Intervals{PollInterval: 1, ReportInterval: 2},
			},
			wantErr: false,
		},
		{
			name: "POLL_INTERVAL=2",
			config: &config.Config{
				Common:    config.CommonConfig{Host: "host"},
				Intervals: config.Intervals{PollInterval: 1, ReportInterval: 1},
			},
			env: map[string]string{"POLL_INTERVAL": "2"},
			want: &config.Config{
				Common:    config.CommonConfig{Host: "host"},
				Intervals: config.Intervals{PollInterval: 2, ReportInterval: 1},
			},
			wantErr: false,
		},
		{
			name: "RATE_LIMIT=5",
			config: &config.Config{
				Common: config.CommonConfig{RateLimit: 1},
			},
			env: map[string]string{"RATE_LIMIT": "5"},
			want: &config.Config{
				Common: config.CommonConfig{RateLimit: 5},
			},
			wantErr: false,
		},
		{
			name: "CRYPTO_KEY=path/to/file.pub",
			config: &config.Config{
				Common: config.CommonConfig{},
			},
			env: map[string]string{"CRYPTO_KEY": "path/to/file.pub"},
			want: &config.Config{
				Common: config.CommonConfig{CryptoKey: "path/to/file.pub"},
			},
			wantErr: false,
		},
		{
			name: "ADDRESS=new-host REPORT_INTERVAL=2 POLL_INTERVAL=2 KEY=abc",
			config: &config.Config{
				Common:    config.CommonConfig{Host: "host"},
				Intervals: config.Intervals{PollInterval: 1, ReportInterval: 1},
			},
			env: map[string]string{"ADDRESS": "new-host", "REPORT_INTERVAL": "2", "POLL_INTERVAL": "2", "KEY": "abc"},
			want: &config.Config{
				Common:    config.CommonConfig{Host: "new-host", Key: "abc"},
				Intervals: config.Intervals{PollInterval: 2, ReportInterval: 2},
			},
			wantErr: false,
		},
	}
	parser := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			for k, v := range tt.env {
				_ = os.Setenv(k, v)
			}
			err := parser.Parse(tt.config)
			assert.Equal(t, tt.want, tt.config)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
