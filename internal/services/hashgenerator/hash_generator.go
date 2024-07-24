package hashgenerator

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type hashGenerator struct {
	key string
}

func New(key string) (*hashGenerator, error) {
	if key == "" {
		return nil, fmt.Errorf("empty key")
	}
	return &hashGenerator{key: key}, nil
}

func (h *hashGenerator) Generate(data []byte) string {
	t := sha256.New()
	t.Write(data)
	return hex.EncodeToString(t.Sum([]byte(h.key)))
}
