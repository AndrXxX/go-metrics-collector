package configpathprovider

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
		name    string
		flags   []string
		want    string
		wantErr bool
	}{
		{
			name:    "Test without path",
			flags:   []string{},
			want:    "",
			wantErr: false,
		},
		{
			name:    "Test with empty path -c=",
			flags:   []string{"-c", ""},
			want:    "",
			wantErr: false,
		},
		{
			name:    "Test with error",
			flags:   []string{"-c"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Test with -c=/path/to/file",
			flags:   []string{"-c", "/path/to/file"},
			want:    "/path/to/file",
			wantErr: false,
		},
		{
			name:    "Test with -config=/path/to/file",
			flags:   []string{"-config", "/path/to/file"},
			want:    "/path/to/file",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			os.Args = os.Args[:1]
			os.Args = append(os.Args[:1], tt.flags...)

			p := FromFlagsProvider{}
			got, err := p.Fetch()
			require.Equal(t, tt.wantErr, err != nil, fmt.Sprintf("%v", err))
			assert.Equal(t, tt.want, got)
		})
	}
}
