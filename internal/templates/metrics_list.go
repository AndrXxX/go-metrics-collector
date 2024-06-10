package templates

const MetricsList = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
	</head>
	<body>
		<h1>{{.Title}}</h1>
		{{ range $metric, $value := .Items }}
			<div><strong>{{ $metric }}:</strong> <span>{{ $value }}</span></div>
		{{else}}
			<div><strong>Список пуст</strong></div>
		{{end}}
	</body>
</html>
`
