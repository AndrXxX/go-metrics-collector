package filters

import (
	"fmt"
	"regexp"

	"golang.org/x/tools/go/analysis"
)

// ByName фильтрует список анализаторов по именам
func ByName(list []*analysis.Analyzer, names []string) ([]*analysis.Analyzer, error) {
	var filtered []*analysis.Analyzer
	for _, check := range list {
		for _, name := range names {
			matched, err := regexp.MatchString(name, check.Name)
			if err != nil {
				return nil, fmt.Errorf("failed to match static check '%s': %v", check.Name, err)
			}
			if matched {
				filtered = append(filtered, check)
			}
		}
	}
	return filtered, nil
}
