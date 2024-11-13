package tlsconfig

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProvider_Fetch(t *testing.T) {
	cert := getPrivateCert()
	tests := []struct {
		name        string
		path        string
		processFile bool
		writeData   []byte
		want        *tls.Config
		wantErr     bool
	}{
		{
			name:    "Test with empty path",
			path:    "",
			want:    nil,
			wantErr: false,
		},
		{
			name:    "Test with not exist file",
			path:    "1.tmp",
			want:    nil,
			wantErr: true,
		},
		{
			name:        "Test with incorrect file",
			path:        "1.tmp",
			processFile: true,
			writeData:   []byte("test"),
			want:        nil,
			wantErr:     true,
		},
		{
			name:        "Test with incorrect file 2",
			path:        "1.tmp",
			processFile: true,
			writeData: func() []byte {
				cert := getPrivateCert()
				var key bytes.Buffer
				_ = pem.Encode(&key, &pem.Block{
					Type:  "RSA PRIVATE KEY",
					Bytes: cert,
				})
				return key.Bytes()
			}(),
			want:    nil,
			wantErr: true,
		},
		{
			name:        "Test with correct file",
			path:        "1.tmp",
			processFile: true,
			writeData:   cert,
			want: func() *tls.Config {
				b, _ := pem.Decode(cert)
				xCert, _ := x509.ParseCertificate(b.Bytes)
				return &tls.Config{
					Certificates: []tls.Certificate{{
						Certificate: [][]byte{xCert.Raw},
						PrivateKey:  xCert,
					}},
				}
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Provider{CryptoKeyPath: tt.path}
			if tt.processFile {
				_ = os.WriteFile(tt.path, tt.writeData, 0644)
			}
			got, err := p.Fetch()
			if tt.processFile {
				_ = os.Remove(tt.path)
			}
			require.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func getPrivateCert() []byte {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization: []string{"Yandex.Praktikum"},
			Country:      []string{"RU"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatal(err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Fatal(err)
	}

	var certPEM bytes.Buffer
	_ = pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	return certPEM.Bytes()
}
