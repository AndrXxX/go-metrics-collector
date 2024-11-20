package flagsparser

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
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
		{
			name:   "-k=abc",
			config: &config.Config{Key: ""},
			flags:  []string{"-k", "abc"},
			want:   &config.Config{Key: "abc"},
		},
		{
			name:   "-crypto-key=/path/to/file.private",
			config: &config.Config{CryptoKey: "some-path"},
			flags:  []string{"-crypto-key", "/path/to/file.private"},
			want:   &config.Config{CryptoKey: "/path/to/file.private"},
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
		err := New().Parse(tt.config)
		assert.Equal(t, tt.want, tt.config)
		assert.NoError(t, err)
	})
}
