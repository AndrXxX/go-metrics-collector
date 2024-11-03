package config

// Config конфигурация multichecker
type Config struct {
	StaticAnalyzers        []string
	ExcludeStaticAnalyzers []string
}
