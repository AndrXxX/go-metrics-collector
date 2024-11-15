package configpathprovider

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testFetcher struct {
	path string
	err  error
}

func (f testFetcher) Fetch() (string, error) {
	return f.path, f.err
}

func TestPathProvider_Fetch(t *testing.T) {
	tests := []struct {
		name     string
		provider PathProvider
		want     string
		wantErr  bool
	}{
		{
			name:     "Test with empty fetchers",
			provider: PathProvider{},
			want:     "",
			wantErr:  false,
		},
		{
			name:     "Test with error on fetch",
			provider: PathProvider{[]Provider{testFetcher{err: fmt.Errorf("test error")}}},
			want:     "",
			wantErr:  true,
		},
		{
			name:     "Test with not empty path",
			provider: PathProvider{[]Provider{testFetcher{path: "/test/path"}}},
			want:     "/test/path",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.provider.Fetch()
			require.Equal(t, tt.wantErr, err != nil, fmt.Sprintf("%v", err))
			assert.Equal(t, tt.want, got)
		})
	}
}
