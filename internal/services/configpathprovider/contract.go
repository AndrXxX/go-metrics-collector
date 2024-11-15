package configpathprovider

// Provider интерфейс получения пути к файлу конфигурации
type Provider interface {
	Fetch() (string, error)
}
