package hashgenerator

import (
	"crypto/sha256"
	"encoding/hex"
)

type hashGenerator struct {
}

func New() *hashGenerator {
	return &hashGenerator{}
}

func (h *hashGenerator) Generate(key string, data []byte) string {
	t := sha256.New()
	t.Write(data)
	return hex.EncodeToString(t.Sum([]byte(key)))
}
