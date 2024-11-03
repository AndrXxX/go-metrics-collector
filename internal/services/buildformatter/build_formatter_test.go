package buildformatter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildFormatter_Format(t *testing.T) {
	tests := []struct {
		name   string
		labels []string
		values []string
		want   []string
	}{
		{
			name:   "Test with Values",
			labels: []string{"Build version", "Build date", "Build commit"},
			values: []string{"1.1", "01.11.2024", "test"},
			want: []string{
				"Build version: 1.1",
				"Build date: 01.11.2024",
				"Build commit: test",
			},
		},
		{
			name:   "Test without Values",
			labels: []string{"Build version", "Build date", "Build commit"},
			values: []string{},
			want: []string{
				"Build version: N/A",
				"Build date: N/A",
				"Build commit: N/A",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := BuildFormatter{
				Labels: tt.labels,
				Values: tt.values,
			}
			assert.Equal(t, tt.want, i.Format())
		})
	}
}

func TestBuildFormatter_combine(t *testing.T) {
	tests := []struct {
		name  string
		label string
		value string
		want  string
	}{
		{
			name:  "Test with empty value",
			label: "Build version",
			want:  "Build version: N/A",
		},
		{
			name:  "Test with value 1.1.1",
			label: "Build version",
			value: "1.1.1",
			want:  "Build version: 1.1.1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := BuildFormatter{}
			assert.Equal(t, tt.want, f.combine(tt.label, tt.value))
		})
	}
}
