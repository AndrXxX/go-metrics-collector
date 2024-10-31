package analyzersprovider

import (
	"testing"

	testifyAnalyzer "github.com/Antonboom/testifylint/analyzer"
	"github.com/kisielk/errcheck/errcheck"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"

	"github.com/AndrXxX/go-metrics-collector/internal/staticlint/config"
	"github.com/AndrXxX/go-metrics-collector/internal/staticlint/services/osexitanalyzer"
	"github.com/AndrXxX/go-metrics-collector/internal/staticlint/vars"
)

func TestAnalyzersProvider_Fetch(t *testing.T) {
	tests := []struct {
		name      string
		config    *config.Config
		wantNames []string
		wantErr   bool
	}{
		{
			name:   "Test with golang analyzers",
			config: &config.Config{},
			wantNames: []string{
				printf.Analyzer.Name,
				shadow.Analyzer.Name,
			},
			wantErr: false,
		},
		{
			name:   "Test with additional analyzers",
			config: &config.Config{},
			wantNames: []string{
				testifyAnalyzer.New().Name,
				errcheck.Analyzer.Name,
				osexitanalyzer.OSExitAnalyzer.Name,
			},
			wantErr: false,
		},
		{
			name: "Test with honnef analyzers",
			config: &config.Config{
				StaticAnalyzers: []string{
					vars.StaticSAAnalyzers,
					vars.StaticSTAnalyzers,
					vars.StaticQFAnalyzers,
				},
			},
			wantNames: []string{
				"SA9008",
				"ST1017",
				"ST1021",
				"QF1005",
			},
			wantErr: false,
		},
		{
			name: "Test with error",
			config: &config.Config{
				StaticAnalyzers: []string{
					"***",
				},
			},
			wantNames: []string{},
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzers, err := AnalyzersProvider{}.Fetch(tt.config)
			require.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				compareAnalyzers(t, analyzers, tt.wantNames)
			}
		})
	}
}

func Test_getAdditionalAnalyzers(t *testing.T) {
	tests := []struct {
		name       string
		checkNames []string
	}{
		{
			name: "Test with testifyAnalyzer, errcheck, osexitanalyzer",
			checkNames: []string{
				"testifylint",
				errcheck.Analyzer.Name,
				osexitanalyzer.OSExitAnalyzer.Name,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compareAnalyzers(t, getAdditionalAnalyzers(), tt.checkNames)
		})
	}
}

func Test_getAnalysisAnalyzers(t *testing.T) {
	tests := []struct {
		name       string
		checkNames []string
	}{
		{
			name: "Test with printf, shadow, shift, structtag",
			checkNames: []string{
				printf.Analyzer.Name,
				shadow.Analyzer.Name,
				shift.Analyzer.Name,
				structtag.Analyzer.Name,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compareAnalyzers(t, getGolangAnalyzers(), tt.checkNames)
		})
	}
}

func compareAnalyzers(t assert.TestingT, list []*analysis.Analyzer, wantNames []string) {
	actualCheckNames := make([]string, len(list))
	for i, analyzer := range list {
		actualCheckNames[i] = analyzer.Name
	}
	for _, name := range wantNames {
		assert.Contains(t, actualCheckNames, name)
	}
}
