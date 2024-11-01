package templates

type BuildData struct {
	Version string
	Date    string
	Commit  string
}

const BuildTemplate = `
Build version: {{if .Version}} "{{.Version}}" {{else}} "N/A" {{end}}
Build date: {{if .Date}} "{{.Date}}" {{else}} "N/A" {{end}}
Build commit: {{if .Commit}} "{{.Commit}}" {{else}} "N/A" {{end}}
`
