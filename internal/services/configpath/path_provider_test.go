package configpath

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

func TestNewProvider(t *testing.T) {
	tests := []struct {
		name string
		opts []func(*pathProvider)
		want *pathProvider
	}{
		{
			name: "Test with empty fetchers",
			opts: []func(*pathProvider){},
			want: &pathProvider{},
		},
		{
			name: "Test with one fetcher",
			opts: []func(*pathProvider){func(p *pathProvider) {
				p.addFetcher(testFetcher{})
			}},
			want: &pathProvider{fetchers: []fetcher{testFetcher{}}},
		},
		{
			name: "Test with two fetchers",
			opts: []func(*pathProvider){
				func(p *pathProvider) { p.addFetcher(testFetcher{}) },
				func(p *pathProvider) { p.addFetcher(testFetcher{}) },
			},
			want: &pathProvider{fetchers: []fetcher{testFetcher{}, testFetcher{}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewProvider(tt.opts...))
		})
	}
}

func TestWithEnv(t *testing.T) {
	tests := []struct {
		name string
		want *pathProvider
	}{
		{
			name: "test OK",
			want: &pathProvider{fetchers: []fetcher{fromEnvProvider{}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewProvider(WithEnv())
			assert.Equal(t, tt.want, p)
		})
	}
}

func TestWithFlags(t *testing.T) {
	tests := []struct {
		name  string
		flags []string
		want  *pathProvider
	}{
		{
			name:  "test with empty flags",
			flags: []string{},
			want:  &pathProvider{fetchers: []fetcher{fromFlagsProvider{flags: []string{}}}},
		},
		{
			name:  "test with flags c and config",
			flags: []string{"c", "config"},
			want:  &pathProvider{fetchers: []fetcher{fromFlagsProvider{flags: []string{"c", "config"}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewProvider(WithFlags(tt.flags...))
			assert.Equal(t, tt.want, p)
		})
	}
}

func Test_pathProvider_Fetch(t *testing.T) {
	tests := []struct {
		name     string
		provider pathProvider
		want     string
		wantErr  bool
	}{
		{
			name:     "Test with empty fetchers",
			provider: pathProvider{},
			want:     "",
			wantErr:  false,
		},
		{
			name:     "Test with error on fetch",
			provider: pathProvider{[]fetcher{testFetcher{err: fmt.Errorf("test error")}}},
			want:     "",
			wantErr:  true,
		},
		{
			name:     "Test with not empty path",
			provider: pathProvider{[]fetcher{testFetcher{path: "/test/path"}}},
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

func Test_pathProvider_addFetcher(t *testing.T) {
	tests := []struct {
		name string
		f    fetcher
		want *pathProvider
	}{
		{
			name: "Test with testFetcher 1",
			f:    testFetcher{path: "test"},
			want: &pathProvider{fetchers: []fetcher{testFetcher{path: "test"}}},
		},
		{
			name: "Test with error testFetcher",
			f:    testFetcher{err: fmt.Errorf("test error")},
			want: &pathProvider{fetchers: []fetcher{testFetcher{err: fmt.Errorf("test error")}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &pathProvider{}
			p.addFetcher(tt.f)
			assert.Equal(t, tt.want, p)
		})
	}
}
