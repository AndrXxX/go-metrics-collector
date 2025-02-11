package envparser

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AndrXxX/go-metrics-collector/internal/staticlint/config"
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
			name: "Empty env",
			config: &config.Config{
				StaticAnalyzers:        []string{"host"},
				ExcludeStaticAnalyzers: []string{},
			},
			env: map[string]string{},
			want: &config.Config{
				StaticAnalyzers:        []string{"host"},
				ExcludeStaticAnalyzers: []string{},
			},
			wantErr: false,
		},
		{
			name: "STATIC_ANALYZERS=SA1.*,SA1032,SA4004 EXCLUDE_STATIC_ANALYZERS=SA1001",
			config: &config.Config{
				StaticAnalyzers: []string{"test"},
			},
			env: map[string]string{"STATIC_ANALYZERS": "SA1.*,SA1032,SA4004", "EXCLUDE_STATIC_ANALYZERS": "SA1001"},
			want: &config.Config{
				StaticAnalyzers: []string{
					"SA1.*",
					"SA1032",
					"SA4004",
				},
				ExcludeStaticAnalyzers: []string{
					"SA1001",
				},
			},
			wantErr: false,
		},
		{
			name: "STATIC_ANALYZERS=",
			config: &config.Config{
				StaticAnalyzers: []string{"test"},
			},
			env: map[string]string{"STATIC_ANALYZERS": ""},
			want: &config.Config{
				StaticAnalyzers: []string{
					"test",
				},
			},
			wantErr: false,
		},
	}
	parser := EnvParser{}
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
