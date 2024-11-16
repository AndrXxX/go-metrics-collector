package configfile

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
)

type testProvider struct {
	path string
	err  error
}

func (p testProvider) Fetch() (string, error) {
	return p.path, p.err
}

func TestConfigFileParser_Parse(t *testing.T) {
	tests := []struct {
		name       string
		c          *config.Config
		provider   pathProvider
		path       string
		fileData   string
		createFile bool
		want       *config.Config
		wantErr    bool
	}{
		{
			name:    "Test without path provider",
			wantErr: false,
		},
		{
			name:     "Test with error while fetch path",
			provider: testProvider{err: fmt.Errorf("test")},
			wantErr:  true,
		},
		{
			name:     "Test with empty path",
			provider: testProvider{path: ""},
			wantErr:  false,
		},
		{
			name:     "Test with not exist file",
			c:        &config.Config{},
			path:     "test.json",
			provider: testProvider{path: "test.json"},
			want:     &config.Config{},
			wantErr:  true,
		},
		{
			name:       "Test with incorrect file",
			c:          &config.Config{},
			path:       "test.json",
			provider:   testProvider{path: "test.json"},
			fileData:   "{\"report_interval\":\"dsfsdf\"}",
			createFile: true,
			want:       &config.Config{},
			wantErr:    true,
		},
		{
			name:       "Test with correct file",
			c:          &config.Config{},
			path:       "test.json",
			provider:   testProvider{path: "test.json"},
			fileData:   "{\"address\":\"localhost:8080\",\"report_interval\":\"1s\",\"poll_interval\":\"1s\",\"crypto_key\":\"/path/to/key.pem\"}",
			createFile: true,
			want: &config.Config{
				Common: config.CommonConfig{
					Host:      "localhost:8080",
					CryptoKey: "/path/to/key.pem",
				},
				Intervals: config.Intervals{
					PollInterval:   1,
					ReportInterval: 1,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.createFile {
				require.NoError(t, os.WriteFile(tt.path, []byte(tt.fileData), 0644))
			}

			p := Parser{tt.provider}
			err := p.Parse(tt.c)

			if tt.createFile {
				_ = os.Remove(tt.path)
			}

			require.Equal(t, tt.wantErr, err != nil, fmt.Errorf("%s", err))
			assert.Equal(t, tt.want, tt.c)
		})
	}
}

func Test_convertDurationToInt(t *testing.T) {
	pointer := func(v int64) *int64 {
		return &v
	}
	tests := []struct {
		name     string
		duration duration
		want     *int64
	}{
		{
			name: "Test nil",
			want: nil,
		},
		{
			name:     "Test with 1s",
			duration: duration{1 * time.Second},
			want:     pointer(1),
		},
		{
			name:     "Test with 100ms",
			duration: duration{100 * time.Millisecond},
			want:     pointer(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, convertDurationToInt(&tt.duration))
		})
	}
}

func Test_set(t *testing.T) {
	type testCase[T comparable] struct {
		name   string
		val    T
		target T
		want   T
	}
	tests := []testCase[any]{
		{
			name:   "Test with val nil",
			val:    nil,
			target: "old_val",
			want:   "old_val",
		},
		{
			name:   "Test with string",
			val:    "new_val",
			target: "old_val",
			want:   "new_val",
		},
		{
			name:   "Test with int",
			val:    10,
			target: 1,
			want:   10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set(&tt.val, &tt.target)
			assert.Equal(t, tt.want, tt.target)
		})
	}
}
