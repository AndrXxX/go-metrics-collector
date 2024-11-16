package configpath

import "fmt"

type pathProvider struct {
	fetchers []fetcher
}

// Fetch метод получения пути к файлу конфигурации
func (p *pathProvider) Fetch() (string, error) {
	for _, fetcher := range p.fetchers {
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

func (p *pathProvider) addFetcher(f fetcher) {
	p.fetchers = append(p.fetchers, f)
}

// WithEnv парсинг пути из переменных среды
func WithEnv() func(*pathProvider) {
	return func(p *pathProvider) {
		p.addFetcher(fromEnvProvider{})
	}
}

// WithFlags парсинг пути из аргументов командной строки
func WithFlags() func(*pathProvider) {
	return func(p *pathProvider) {
		p.addFetcher(fromFlagsProvider{})
	}
}

// NewProvider возвращает сервис для получения пути к файлу конфигурации
func NewProvider(opts ...func(*pathProvider)) *pathProvider {
	p := &pathProvider{}
	for _, opt := range opts {
		opt(p)
	}
	return p
}
