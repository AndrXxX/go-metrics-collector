package hashgenerator

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
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

func (h *hashGenerator) Generate(data []byte) ([]byte, error) {
	//key, err := generateRandom(2 * aes.BlockSize) // AES-256 (32 байта)
	//if err != nil {
	//	return "", fmt.Errorf("error on generate random %w", err)
	//}

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
