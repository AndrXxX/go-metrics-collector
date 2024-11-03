package buildformatter

import (
	"fmt"
)

const emptyInfoTemplate = "N/A"

// BuildFormatter сервис для форматирования информации о текущей сборке
type BuildFormatter struct {
	Labels []string
	Values []string
}

// Format форматирует информацию о текущей сборке в строки
func (f BuildFormatter) Format() []string {
	var result []string
	for i, label := range f.Labels {
		value := ""
		if len(f.Values) >= i+1 {
			value = f.Values[i]
		}
		result = append(result, f.combine(label, value))
	}
	return result
}

func (f BuildFormatter) combine(label string, value string) string {
	if value == "" {
		value = emptyInfoTemplate
	}
	return fmt.Sprintf("%s: %s", label, value)
}
