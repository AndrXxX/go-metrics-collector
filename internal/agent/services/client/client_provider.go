package client

import (
	"fmt"
	"net/http"
)

// Provider сервис для получения http.client с учетом ассиметричного шифрования
type Provider struct {
	ConfProvider tlsConfigProvider
}

// Fetch возвращает http.Client с учетом ассиметричного шифрования
func (p Provider) Fetch() (*http.Client, error) {
	if p.ConfProvider == nil {
		return &http.Client{}, nil
	}
	conf, err := p.ConfProvider.Fetch()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tls configuration: %w", err)
	}
	if conf == nil {
		return &http.Client{}, nil
	}
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: conf,
		},
	}, nil
}
