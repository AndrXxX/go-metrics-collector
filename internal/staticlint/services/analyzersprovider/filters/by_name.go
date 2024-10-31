package filters

import (
	"fmt"
	"regexp"

	"golang.org/x/tools/go/analysis"
)

// ByName фильтрует список анализаторов по именам
func ByName(list []*analysis.Analyzer, names []string, excludeNames []string) ([]*analysis.Analyzer, error) {
	var filtered []*analysis.Analyzer
	for _, analyzer := range list {
		matched, err := granted(analyzer, names, excludeNames)
		if err != nil {
			return nil, err
		}
		if matched {
			filtered = append(filtered, analyzer)
		}
	}
	return filtered, nil
}

func granted(analyzer *analysis.Analyzer, names []string, excludeNames []string) (bool, error) {
	for _, name := range excludeNames {
		result, err := regexp.MatchString(name, analyzer.Name)
		if err != nil {
			return false, fmt.Errorf("failed to match static analyzer '%s': %v", analyzer.Name, err)
		}
		if result {
			return false, nil
		}
	}
	for _, name := range names {
		result, err := regexp.MatchString(name, analyzer.Name)
		if err != nil {
			return false, fmt.Errorf("failed to match static analyzer '%s': %v", analyzer.Name, err)
		}
		if result {
			return true, nil
		}
	}
	return false, nil
}
