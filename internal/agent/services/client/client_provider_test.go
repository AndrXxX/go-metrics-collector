package client

import (
	"crypto/tls"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testConfProvider struct {
	err error
	c   *tls.Config
}

func (p testConfProvider) Fetch() (*tls.Config, error) {
	return p.c, p.err
}

func TestProvider_Fetch(t *testing.T) {
	tests := []struct {
		name         string
		ConfProvider tlsConfigProvider
		want         *http.Client
		wantErr      bool
	}{
		{
			name:    "Test without ConfProvider",
			want:    &http.Client{},
			wantErr: false,
		},
		{
			name:         "Test with error while fetch config",
			ConfProvider: testConfProvider{err: errors.New("test")},
			want:         nil,
			wantErr:      true,
		},
		{
			name:         "Test with nil config",
			ConfProvider: testConfProvider{},
			want:         &http.Client{},
			wantErr:      false,
		},
		{
			name:         "Test with right config",
			ConfProvider: testConfProvider{c: &tls.Config{InsecureSkipVerify: true}},
			want: &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Provider{ConfProvider: tt.ConfProvider}
			got, err := p.Fetch()
			require.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
