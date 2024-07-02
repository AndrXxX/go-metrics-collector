package main

import (
	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_parseEnv(t *testing.T) {
	tests := []struct {
		name   string
		config *config.Config
		env    map[string]string
		want   *config.Config
	}{
		{
			name:   "Empty env",
			config: &config.Config{Host: "host"},
			env:    map[string]string{},
			want:   &config.Config{Host: "host"},
		},
		{
			name:   "ADDRESS=new-host",
			config: &config.Config{Host: "host"},
			env:    map[string]string{"ADDRESS": "new-host"},
			want:   &config.Config{Host: "new-host"},
		},
		{
			name:   "STORE_INTERVAL=5",
			config: &config.Config{StoreInterval: 1},
			env:    map[string]string{"STORE_INTERVAL": "5"},
			want:   &config.Config{StoreInterval: 5},
		},
		{
			name:   "FILE_STORAGE_PATH=/tmp/metrics-test-db.json",
			config: &config.Config{FileStoragePath: "/tmp/temp"},
			env:    map[string]string{"FILE_STORAGE_PATH": "/tmp/metrics-test-db.json"},
			want:   &config.Config{FileStoragePath: "/tmp/metrics-test-db.json"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			for k, v := range tt.env {
				_ = os.Setenv(k, v)
			}
			parseEnv(tt.config)
			assert.Equal(t, tt.want, tt.config)
		})
	}
}
