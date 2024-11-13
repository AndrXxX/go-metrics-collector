package tlsconfig

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
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
	b, _ := pem.Decode(file)
	if b == nil {
		return nil, fmt.Errorf("failed to decode crypto key")
	}
	cert, err := x509.ParseCertificate(b.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{cert.Raw},
			PrivateKey:  cert,
		}},
	}, nil
}
