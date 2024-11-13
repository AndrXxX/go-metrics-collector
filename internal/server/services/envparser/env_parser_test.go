package envparser

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
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
			name:    "Empty env",
			config:  &config.Config{Host: "host"},
			env:     map[string]string{},
			want:    &config.Config{Host: "host"},
			wantErr: false,
		},
		{
			name:    "ADDRESS=new-host",
			config:  &config.Config{Host: "host"},
			env:     map[string]string{"ADDRESS": "new-host"},
			want:    &config.Config{Host: "new-host"},
			wantErr: false,
		},
		{
			name:    "STORE_INTERVAL=5",
			config:  &config.Config{StoreInterval: 1},
			env:     map[string]string{"STORE_INTERVAL": "5"},
			want:    &config.Config{StoreInterval: 5},
			wantErr: false,
		},
		{
			name:    "FILE_STORAGE_PATH=/tmp/metrics-test-db.json",
			config:  &config.Config{FileStoragePath: "/tmp/temp"},
			env:     map[string]string{"FILE_STORAGE_PATH": "/tmp/metrics-test-db.json"},
			want:    &config.Config{FileStoragePath: "/tmp/metrics-test-db.json"},
			wantErr: false,
		},
		{
			name:    "RESTORE=1",
			config:  &config.Config{Restore: false},
			env:     map[string]string{"RESTORE": "1"},
			want:    &config.Config{Restore: true},
			wantErr: false,
		},
		{
			name:    "DATABASE_DSN=test",
			config:  &config.Config{DatabaseDSN: ""},
			env:     map[string]string{"DATABASE_DSN": "test"},
			want:    &config.Config{DatabaseDSN: "test"},
			wantErr: false,
		},
		{
			name:    "KEY=abc",
			config:  &config.Config{Key: ""},
			env:     map[string]string{"KEY": "abc"},
			want:    &config.Config{Key: "abc"},
			wantErr: false,
		},
		{
			name:    "CRYPTO_KEY=path/to/file.pub",
			config:  &config.Config{CryptoKey: "some-path"},
			env:     map[string]string{"CRYPTO_KEY": "path/to/file.pub"},
			want:    &config.Config{CryptoKey: "path/to/file.pub"},
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
