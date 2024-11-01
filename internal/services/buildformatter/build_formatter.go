package buildformatter

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/AndrXxX/go-metrics-collector/internal/services/buildformatter/templates"
)

// BuildFormatter сервис для форматирования информации о текущей сборке
type BuildFormatter struct {
	Version string
	Date    string
	Commit  string
	Buffer  buffer
}

// Format форматирует информацию о текущей сборке в строку
func (i BuildFormatter) Format() (string, error) {
	if i.Buffer == nil {
		i.Buffer = &bytes.Buffer{}
	}
	var tmpl = template.Must(template.New("build").Parse(templates.BuildTemplate))
	err := tmpl.Execute(i.Buffer, templates.BuildData{
		Version: i.Version,
		Date:    i.Date,
		Commit:  i.Commit,
	})
	if err != nil {
		return "", fmt.Errorf("template execution error: %w", err)
	}
	return i.Buffer.String(), nil
}
