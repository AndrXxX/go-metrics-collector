package hashgenerator

import (
	"crypto/sha256"
	"encoding/hex"
)

type sha256Generator struct {
}

func (h *sha256Generator) Generate(key string, data []byte) string {
	t := sha256.New()
	t.Write(data)
	return hex.EncodeToString(t.Sum([]byte(key)))
}
