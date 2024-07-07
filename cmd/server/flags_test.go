package main

import (
	"flag"
	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
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
			name:   "Empty flags",
			config: &config.Config{Host: "host"},
			flags:  []string{},
			want:   &config.Config{Host: "host"},
		},
		{
			name:   "-a=new-host",
			config: &config.Config{Host: "host"},
			flags:  []string{"-a", "new-host"},
			want:   &config.Config{Host: "new-host"},
		},
		{
			name:   "-i=5",
			config: &config.Config{StoreInterval: 1},
			flags:  []string{"-i", "5"},
			want:   &config.Config{StoreInterval: 5},
		},
		{
			name:   "-f=/tmp/1.j",
			config: &config.Config{FileStoragePath: "/tmp/2.j"},
			flags:  []string{"-f", "/tmp/1.j"},
			want:   &config.Config{FileStoragePath: "/tmp/1.j"},
		},
		{
			name:   "-r=0",
			config: &config.Config{Restore: false},
			flags:  []string{"-r", "0"},
			want:   &config.Config{Restore: true},
		},
		{
			name:   "-d=test",
			config: &config.Config{DatabaseDSN: ""},
			flags:  []string{"-d", "test"},
			want:   &config.Config{DatabaseDSN: "test"},
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
