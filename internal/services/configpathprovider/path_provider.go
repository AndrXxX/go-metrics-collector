package configpathprovider

import "fmt"

// PathProvider сервис для получения пути к файлу конфигурации, использующий переданные Provider
type PathProvider struct {
	Fetchers []Provider
}

// Fetch метод получения пути к файлу конфигурации
func (p PathProvider) Fetch() (string, error) {
	for _, fetcher := range p.Fetchers {
		path, err := fetcher.Fetch()
		if err != nil {
			return "", fmt.Errorf("failed to fetch config path: %w", err)
		}
		if path != "" {
			return path, nil
		}
	}
	return "", nil
}
