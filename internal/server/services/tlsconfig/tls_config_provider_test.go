package tlsconfig

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProvider_Fetch(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		processFile bool
		writeData   []byte
		want        *tls.Config
		wantErr     bool
	}{
		{
			name:    "Test with empty path",
			path:    "",
			want:    nil,
			wantErr: false,
		},
		{
			name:    "Test with not exist file",
			path:    "1.tmp",
			want:    nil,
			wantErr: true,
		},
		{
			name:        "Test with correct file",
			path:        "1.tmp",
			processFile: true,
			writeData:   []byte("test"),
			want: &tls.Config{
				ClientCAs: func() *x509.CertPool {
					certPool := x509.NewCertPool()
					certPool.AppendCertsFromPEM([]byte("test"))
					return certPool
				}(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Provider{CryptoKeyPath: tt.path}
			if tt.processFile {
				_ = os.WriteFile(tt.path, tt.writeData, 0644)
			}
			got, err := p.Fetch()
			if tt.processFile {
				_ = os.Remove(tt.path)
			}
			require.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
