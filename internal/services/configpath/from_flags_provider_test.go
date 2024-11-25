package configpath

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlagPathProvider_Fetch(t *testing.T) {
	tests := []struct {
		name      string
		flagNames []string
		flags     []string
		want      string
		wantErr   bool
	}{
		{
			name:      "Test with empty flagNames",
			flagNames: make([]string, 0),
			flags:     []string{},
			want:      "",
			wantErr:   false,
		},
		{
			name:      "Test without path",
			flagNames: []string{"c"},
			flags:     []string{},
			want:      "",
			wantErr:   false,
		},
		{
			name:      "Test with empty path -c=",
			flagNames: []string{"c"},
			flags:     []string{"-c", ""},
			want:      "",
			wantErr:   false,
		},
		{
			name:      "Test with error",
			flagNames: []string{"c"},
			flags:     []string{"-c"},
			want:      "",
			wantErr:   true,
		},
		{
			name:      "Test with -c=/path/to/file",
			flagNames: []string{"c"},
			flags:     []string{"-c", "/path/to/file"},
			want:      "/path/to/file",
			wantErr:   false,
		},
		{
			name:      "Test with -config=/path/to/file",
			flagNames: []string{"c", "config"},
			flags:     []string{"-config", "/path/to/file"},
			want:      "/path/to/file",
			wantErr:   false,
		},
		{
			name:      "Test with -config=/path/to/file and additional flags",
			flagNames: []string{"c", "config"},
			flags:     []string{"-config", "/path/to/file", "-a", "someValue"},
			want:      "/path/to/file",
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			os.Args = os.Args[:1]
			os.Args = append(os.Args[:1], tt.flags...)

			p := fromFlagsProvider{tt.flagNames}
			got, err := p.Fetch()
			require.Equal(t, tt.wantErr, err != nil, fmt.Sprintf("%v", err))
			assert.Equal(t, tt.want, got)
		})
	}
}
