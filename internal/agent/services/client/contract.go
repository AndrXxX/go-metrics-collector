package client

import "crypto/tls"

type tlsConfigProvider interface {
	Fetch() (*tls.Config, error)
}
