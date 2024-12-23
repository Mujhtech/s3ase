package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

type AesGcm struct {
	block   cipher.Block
	compact bool
}

func NewAesGcm(key string) (Encrypt, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("key must be 32 bytes")
	}
	keyByte := []byte(key)

	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return nil, err
	}
	return &AesGcm{
		block:   block,
		compact: true,
	}, nil
}

func (e *AesGcm) Encrypt(plaintext []byte) (string, error) {
	gcm, err := cipher.NewGCM(e.block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	return string(gcm.Seal(nonce, nonce, plaintext, nil)), nil
}

func (e *AesGcm) Decrypt(ciphertext string) (string, error) {

	ciphertextByte := []byte(ciphertext)

	gcm, err := cipher.NewGCM(e.block)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < gcm.NonceSize() {

		if e.compact {
			return string(ciphertext), nil
		}
		return "", fmt.Errorf("ciphertext too short")
	}

	plaintext, err := gcm.Open(nil,
		ciphertextByte[:gcm.NonceSize()],
		ciphertextByte[gcm.NonceSize():],
		nil,
	)

	if err != nil && e.compact {
		return string(ciphertext), nil
	}
	return string(plaintext), err
}
