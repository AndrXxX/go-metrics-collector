package hashgenerator

import (
	"crypto/aes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"slices"
)

var keySizes = []int{aes.BlockSize, 1.5 * aes.BlockSize, 2 * aes.BlockSize}

type hashGenerator struct {
	key string
}

func New(key string) (*hashGenerator, error) {
	if !slices.Contains(keySizes, len(key)) {
		return nil, fmt.Errorf("invalid key size %s", key)
	}
	return &hashGenerator{key: key}, nil
}

func (h *hashGenerator) Generate(data []byte) string {
	t := sha256.New()
	t.Write(data)
	return hex.EncodeToString(t.Sum([]byte(h.key)))
}
