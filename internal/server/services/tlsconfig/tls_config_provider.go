package tlsconfig

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

type Provider struct {
	CryptoKeyPath string
}

func (p Provider) Fetch() (*tls.Config, error) {
	if p.CryptoKeyPath == "" {
		return nil, nil
	}
	file, err := os.ReadFile(p.CryptoKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read crypto key file: %w", err)
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(file)
	return &tls.Config{
		ClientCAs: certPool,
	}, nil
}
