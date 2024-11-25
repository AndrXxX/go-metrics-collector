package configpath

import (
	fl "flag"
	"fmt"
	"os"
	"slices"
	"strings"
)

type fromFlagsProvider struct {
	flags []string
}

// Fetch метод получения пути к файлу конфигурации
func (p fromFlagsProvider) Fetch() (string, error) {
	fs := fl.NewFlagSet("path", fl.ContinueOnError)
	paths := make([]*string, len(p.flags))
	flags := make([]string, len(p.flags))
	for i := range p.flags {
		paths[i] = fs.String(p.flags[i], "", fmt.Sprintf("Path to config JSON file -%s", p.flags[i]))
		flags[i] = "-" + p.flags[i]
	}
	args := make([]string, 0)
	rawArgs := os.Args[1:]
	for i := 0; i < len(rawArgs); i += 2 {
		if slices.Contains(flags, rawArgs[i]) {
			args = append(args, rawArgs[i])
			if i+1 < len(rawArgs) && !strings.Contains(rawArgs[i+1], "-") {
				args = append(args, rawArgs[i+1])
			}
		}
	}
	if err := fs.Parse(args); err != nil {
		return "", fmt.Errorf("failed to parse flag: %s", err)
	}
	for _, path := range paths {
		if *path != "" {
			return *path, nil
		}
	}
	return "", nil
}
