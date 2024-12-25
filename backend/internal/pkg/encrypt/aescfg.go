package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

type AesCfb struct {
	block cipher.Block
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
	return &AesCfb{block: block}, nil
}

func (e *AesCfb) Encrypt(plainText []byte) (string, error) {
	const maxSize = 64 * 1024 * 1024 // 64 MB
	if len(plainText) > maxSize {
		return "", fmt.Errorf("plainText too large")
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	encrypter := cipher.NewCFBEncrypter(e.block, iv)
	encrypter.XORKeyStream(cipherText[aes.BlockSize:], []byte(plainText))

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func (e *AesCfb) Decrypt(cipherText string) (string, error) {

	cipherTextBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	iv := cipherTextBytes[:aes.BlockSize]
	cipherTextBytes = cipherTextBytes[aes.BlockSize:]

	decrypter := cipher.NewCFBDecrypter(e.block, iv)
	decrypter.XORKeyStream(cipherTextBytes, cipherTextBytes)

	return string(cipherTextBytes), nil
}
