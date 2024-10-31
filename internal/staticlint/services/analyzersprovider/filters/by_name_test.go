package filters

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis"
)

func TestByName(t *testing.T) {
	tests := []struct {
		name      string
		analyzers []*analysis.Analyzer
		names     []string
		want      []string
		wantErr   bool
	}{
		{
			name: "filter with test_name",
			analyzers: []*analysis.Analyzer{
				{Name: "test_name"},
				{Name: "test_name2"},
			},
			names: []string{"test_name"},
			want: []string{
				"test_name",
				"test_name2",
			},
			wantErr: false,
		},
		{
			name: "filter with ST.* and SA.*",
			analyzers: []*analysis.Analyzer{
				{Name: "ST1005"},
				{Name: "ST3004"},
				{Name: "SA3423"},
				{Name: "SB3423"},
			},
			names: []string{"ST.*", "SA.*"},
			want: []string{
				"ST1005",
				"ST3004",
				"SA3423",
			},
			wantErr: false,
		},
		{
			name: "test with error",
			analyzers: []*analysis.Analyzer{
				{Name: "test_name"},
			},
			names:   []string{"***"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzers, err := ByName(tt.analyzers, tt.names)
			require.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				assert.Equal(t, len(tt.want), len(analyzers))
				for _, a := range analyzers {
					assert.Contains(t, tt.want, a.Name)
				}
			}
		})
	}
}
