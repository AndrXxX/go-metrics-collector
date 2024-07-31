package main

import (
	"flag"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type testCase struct {
	name   string
	config *config.Config
	flags  []string
	want   *config.Config
}

func Test_parseFlags(t *testing.T) {
	tests := []testCase{
		{
			name: "Empty flags",
			config: &config.Config{
				Common:    config.CommonConfig{Host: "host"},
				Intervals: config.Intervals{PollInterval: 1, ReportInterval: 1},
			},
			flags: []string{},
			want: &config.Config{
				Common:    config.CommonConfig{Host: "host"},
				Intervals: config.Intervals{PollInterval: 1, ReportInterval: 1},
			},
		},
		{
			name: "-a=new-host",
			config: &config.Config{
				Common:    config.CommonConfig{Host: "host"},
				Intervals: config.Intervals{PollInterval: 1, ReportInterval: 1},
			},
			flags: []string{"-a", "new-host"},
			want: &config.Config{
				Common:    config.CommonConfig{Host: "new-host"},
				Intervals: config.Intervals{PollInterval: 1, ReportInterval: 1},
			},
		},
		{
			name: "-r=2",
			config: &config.Config{
				Common:    config.CommonConfig{Host: "host"},
				Intervals: config.Intervals{PollInterval: 1, ReportInterval: 1},
			},
			flags: []string{"-r", "2"},
			want: &config.Config{
				Common:    config.CommonConfig{Host: "host"},
				Intervals: config.Intervals{PollInterval: 1, ReportInterval: 2},
			},
		},
		{
			name: "-p=2",
			config: &config.Config{
				Common:    config.CommonConfig{Host: "host"},
				Intervals: config.Intervals{PollInterval: 1, ReportInterval: 1},
			},
			flags: []string{"-p", "2"},
			want: &config.Config{
				Common:    config.CommonConfig{Host: "host"},
				Intervals: config.Intervals{PollInterval: 2, ReportInterval: 1},
			},
		},
		{
			name: "-l=2",
			config: &config.Config{
				Common: config.CommonConfig{RateLimit: 1},
			},
			flags: []string{"-l", "2"},
			want: &config.Config{
				Common: config.CommonConfig{RateLimit: 2},
			},
		},
		{
			name: "-a=new-host -r=2 -p=2 -k=abc",
			config: &config.Config{
				Common:    config.CommonConfig{Host: "host"},
				Intervals: config.Intervals{PollInterval: 1, ReportInterval: 1},
			},
			flags: []string{"-a", "new-host", "-p", "2", "-r", "2", "-k", "abc"},
			want: &config.Config{
				Common:    config.CommonConfig{Host: "new-host", Key: "abc"},
				Intervals: config.Intervals{PollInterval: 2, ReportInterval: 2},
			},
		},
	}
	for _, tt := range tests {
		run(t, tt)
	}
}

func run(t *testing.T, tt testCase) {
	t.Run(tt.name, func(t *testing.T) {
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		os.Args = os.Args[:1]
		os.Args = append(os.Args[:1], tt.flags...)
		parseFlags(tt.config)
		assert.Equal(t, tt.want, tt.config)
	})
}
