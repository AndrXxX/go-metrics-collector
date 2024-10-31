package honnef

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis"
	"honnef.co/go/tools/analysis/lint"

	"github.com/AndrXxX/go-metrics-collector/internal/staticlint/config"
	"github.com/AndrXxX/go-metrics-collector/internal/staticlint/vars"
)

func TestAnalyzers(t *testing.T) {
	tests := []struct {
		name       string
		checkNames []string
		config     *config.Config
		wantErr    bool
	}{
		{
			name: "Test with defaults",
			checkNames: []string{
				"SA1000",
				"SA1032",
				"SA4004",
				"SA9008",
				"ST1017",
				"ST1021",
				"QF1005",
				"QF1011",
			},
			config: &config.Config{StaticChecks: []string{
				vars.StaticSAChecks,
				vars.StaticSTChecks,
				vars.StaticQFChecks,
			}},
			wantErr: false,
		},
		{
			name:       "Test with error",
			checkNames: []string{},
			config: &config.Config{StaticChecks: []string{
				"***",
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checks, err := Analyzers(tt.config)
			require.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				compareChecks(t, checks, tt.checkNames)
			}
		})
	}
}

func compareChecks(t assert.TestingT, checks []*analysis.Analyzer, wantChecks []string) {
	actualCheckNames := make([]string, len(checks))
	for i, check := range checks {
		actualCheckNames[i] = check.Name
	}
	for _, checkName := range wantChecks {
		assert.Contains(t, actualCheckNames, checkName)
	}
}

func Test_convert(t *testing.T) {
	tests := []struct {
		name string
		list []*lint.Analyzer
		want []*analysis.Analyzer
	}{
		{
			name: "Test with empty slice",
			list: []*lint.Analyzer{},
			want: []*analysis.Analyzer{},
		},
		{
			name: "Test with test and test2",
			list: []*lint.Analyzer{
				{Analyzer: &analysis.Analyzer{Name: "test"}},
				{Analyzer: &analysis.Analyzer{Name: "test2"}},
			},
			want: []*analysis.Analyzer{
				{Name: "test"},
				{Name: "test2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, convert(tt.list))
		})
	}
}
