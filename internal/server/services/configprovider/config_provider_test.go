package configprovider

import (
	"errors"
	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

type tempParser struct {
	err  error
	host string
}

func (p tempParser) Parse(c *config.Config) error {
	c.Host = p.host
	return p.err
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		parsers []parser
		want    *configProvider
	}{
		{
			name:    "Test with empty parsers",
			parsers: []parser{},
			want:    &configProvider{parsers: []parser{}},
		},
		{
			name:    "Test with temp parser",
			parsers: []parser{tempParser{}},
			want:    &configProvider{parsers: []parser{tempParser{}}},
		},
		{
			name:    "Test with two temp parsers",
			parsers: []parser{tempParser{}, tempParser{}},
			want:    &configProvider{parsers: []parser{tempParser{}, tempParser{}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, New(tt.parsers...))
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		parsers []parser
		wantErr bool
	}{
		{
			name:    "Test with err parser",
			parsers: []parser{tempParser{err: errors.New("err")}},
			wantErr: true,
		},
		{
			name:    "Test with no err parser",
			parsers: []parser{tempParser{}},
			wantErr: false,
		},
		{
			name:    "Test with validate err",
			parsers: []parser{tempParser{host: "-"}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := New(tt.parsers...)
			_, err := provider.Fetch()
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
