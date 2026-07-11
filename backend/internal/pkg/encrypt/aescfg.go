package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
)

const authenticatedCipherPrefix = "gcm:"

type AES struct {
	block cipher.Block
	aead  cipher.AEAD
}

func NewAesCfb(key string) (Encrypt, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("key must be 32 bytes")
	}
	b := []byte(key)
	block, err := aes.NewCipher(b)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &AES{block: block, aead: aead}, nil
}

func (e *AES) Encrypt(plainText []byte) (string, error) {
	const maxSize = 64 * 1024 * 1024 // 64 MB
	if len(plainText) > maxSize {
		return "", fmt.Errorf("plainText too large")
	}

	nonce := make([]byte, e.aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	cipherText := e.aead.Seal(nonce, nonce, plainText, nil)
	return authenticatedCipherPrefix + base64.StdEncoding.EncodeToString(cipherText), nil
}

func (e *AES) Decrypt(cipherText string) (string, error) {
	if !strings.HasPrefix(cipherText, authenticatedCipherPrefix) {
		return e.decryptLegacyCFB(cipherText)
	}
	cipherTextBytes, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(cipherText, authenticatedCipherPrefix))
	if err != nil {
		return "", err
	}
	if len(cipherTextBytes) < e.aead.NonceSize()+e.aead.Overhead() {
		return "", fmt.Errorf("ciphertext too short")
	}
	nonce, encrypted := cipherTextBytes[:e.aead.NonceSize()], cipherTextBytes[e.aead.NonceSize():]
	plainText, err := e.aead.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return "", fmt.Errorf("failed to authenticate ciphertext: %w", err)
	}
	return string(plainText), nil
}
