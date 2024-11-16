package configpath

import (
	fl "flag"
	"fmt"
	"os"
)

type fromFlagsProvider struct {
}

// Fetch метод получения пути к файлу конфигурации
func (p fromFlagsProvider) Fetch() (string, error) {
	fs := fl.NewFlagSet("path", fl.ContinueOnError)
	path := fs.String("c", "", "Path to config JSON file")
	path2 := fs.String("config", "", "Path to config JSON file")
	if err := fs.Parse(os.Args[1:]); err != nil {
		return "", fmt.Errorf("failed to parse flag: %s", err)
	}
	if *path == "" {
		path = path2
	}
	return *path, nil
}
