package flagsparser

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AndrXxX/go-metrics-collector/internal/staticlint/config"
)

type testCase struct {
	name   string
	config *config.Config
	flags  []string
	want   *config.Config
}

func Test_parseFlags(t *testing.T) {
	tests := []testCase{
		{
			name: "Empty flags",
			config: &config.Config{
				StaticAnalyzers: []string{"test"},
			},
			flags: []string{},
			want: &config.Config{
				StaticAnalyzers: []string{"test"},
			},
		},
		{
			name: "-sa=SA1000,SA1032,SA4004",
			config: &config.Config{
				StaticAnalyzers: []string{"test"},
			},
			flags: []string{"-sa", "SA1000,SA1032,SA4004"},
			want: &config.Config{
				StaticAnalyzers: []string{
					"SA1000",
					"SA1032",
					"SA4004",
				},
			},
		},
	}
	for _, tt := range tests {
		run(t, tt)
	}
}

func run(t *testing.T, tt testCase) {
	t.Run(tt.name, func(t *testing.T) {
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		os.Args = os.Args[:1]
		os.Args = append(os.Args[:1], tt.flags...)
		err := FlagsParser{}.Parse(tt.config)
		assert.Equal(t, tt.want, tt.config)
		assert.NoError(t, err)
	})
}

func Example_flagsParser_Parse() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = os.Args[:1]
	os.Args = append(os.Args[:1], []string{"-sa", "SA1000,SA1032,SA4004"}...)

	c := config.NewConfig()
	_ = FlagsParser{}.Parse(c)

	fmt.Println(c.StaticAnalyzers)

	// Output:
	// [SA1000 SA1032 SA4004]
}
