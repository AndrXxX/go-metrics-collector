package configpath

import (
	fl "flag"
	"fmt"
	"os"
)

type fromFlagsProvider struct {
	flags []string
}

// Fetch метод получения пути к файлу конфигурации
func (p fromFlagsProvider) Fetch() (string, error) {
	fs := fl.NewFlagSet("path", fl.ContinueOnError)
	paths := make([]*string, len(p.flags))
	for i := range p.flags {
		paths[i] = fs.String(p.flags[i], "", fmt.Sprintf("Path to config JSON file -%s", p.flags[i]))
	}
	if err := fs.Parse(os.Args[1:]); err != nil {
		return "", fmt.Errorf("failed to parse flag: %s", err)
	}
	for _, path := range paths {
		if *path != "" {
			return *path, nil
		}
	}
	return "", nil
}
