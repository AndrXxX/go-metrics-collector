package main

import (
	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
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
			name: "ADDRESS=new-host REPORT_INTERVAL=2 POLL_INTERVAL=2",
			config: &config.Config{
				Common:    config.CommonConfig{Host: "host"},
				Intervals: config.Intervals{PollInterval: 1, ReportInterval: 1},
			},
			env: map[string]string{"ADDRESS": "new-host", "REPORT_INTERVAL": "2", "POLL_INTERVAL": "2"},
			want: &config.Config{
				Common:    config.CommonConfig{Host: "new-host"},
				Intervals: config.Intervals{PollInterval: 2, ReportInterval: 2},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			for k, v := range tt.env {
				_ = os.Setenv(k, v)
			}
			err := parseEnv(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, tt.config)
		})
	}
}
