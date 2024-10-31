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
			config: &config.Config{StaticAnalyzers: []string{
				vars.StaticSAAnalyzers,
				vars.StaticSTAnalyzers,
				vars.StaticQFAnalyzers,
			}},
			wantErr: false,
		},
		{
			name:       "Test with error",
			checkNames: []string{},
			config: &config.Config{StaticAnalyzers: []string{
				"***",
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := Analyzers(tt.config)
			require.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				compareAnalyzers(t, list, tt.checkNames)
			}
		})
	}
}

func compareAnalyzers(t assert.TestingT, list []*analysis.Analyzer, wantNames []string) {
	actualNames := make([]string, len(list))
	for i, analyzer := range list {
		actualNames[i] = analyzer.Name
	}
	for _, name := range wantNames {
		assert.Contains(t, actualNames, name)
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
