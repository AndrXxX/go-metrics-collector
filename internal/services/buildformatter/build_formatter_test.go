package buildformatter

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testBuffer struct {
}

func (b testBuffer) Write(_ []byte) (n int, err error) {
	return 0, fmt.Errorf("error")
}

func (b testBuffer) String() string {
	return ""
}

func TestBuildFormatter_Format(t *testing.T) {
	type fields struct {
		Version string
		Date    string
		Commit  string
		Buffer  buffer
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name:   "Test with values",
			fields: fields{"1.1", "01.11.2024", "test", nil},
			want: `
Build version:  "1.1" 
Build date:  "01.11.2024" 
Build commit:  "test" 
`,
			wantErr: false,
		},
		{
			name:   "Test without values",
			fields: fields{},
			want: `
Build version:  "N/A" 
Build date:  "N/A" 
Build commit:  "N/A" 
`,
			wantErr: false,
		},
		{
			name:    "Test with error",
			fields:  fields{Buffer: testBuffer{}},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := BuildFormatter{
				Version: tt.fields.Version,
				Date:    tt.fields.Date,
				Commit:  tt.fields.Commit,
				Buffer:  tt.fields.Buffer,
			}
			str, err := i.Format()
			require.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, str)
		})
	}
}
