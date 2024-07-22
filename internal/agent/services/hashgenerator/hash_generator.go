package hashgenerator

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

type hashGenerator struct {
	key string
}

func New(key string) *hashGenerator {
	return &hashGenerator{key: key}
}

func (h *hashGenerator) Generate(data []byte) ([]byte, error) {
	//key, err := generateRandom(2 * aes.BlockSize) // AES-256 (32 байта)
	//if err != nil {
	//	return "", fmt.Errorf("error on generate random %w", err)
	//}
	if h.key == "" {
		return nil, fmt.Errorf("empty key provided")
	}

	aesBlock, err := aes.NewCipher([]byte(h.key))
	if err != nil {
		return nil, fmt.Errorf("error on create aes block %w", err)
	}

	aesGCM, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return nil, fmt.Errorf("error on create aes gcm %w", err)
	}

	nonce, err := generateRandom(aesGCM.NonceSize())
	if err != nil {
		return nil, fmt.Errorf("error on create nonce %w", err)
	}

	dst := aesGCM.Seal(nil, nonce, data, nil)
	return dst, nil
}

func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
