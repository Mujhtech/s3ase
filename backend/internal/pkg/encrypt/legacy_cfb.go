package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

// decryptLegacyCFB is retained only to read values written before AES-GCM was
// introduced. All new ciphertext is authenticated and versioned by Encrypt.
func (e *AES) decryptLegacyCFB(encoded string) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	if len(cipherText) < aes.BlockSize {
		return "", fmt.Errorf("legacy ciphertext too short")
	}
	iv, encrypted := cipherText[:aes.BlockSize], cipherText[aes.BlockSize:]
	//nolint:staticcheck // Required for backward-compatible decryption only; encryption always uses AES-GCM.
	decrypter := cipher.NewCFBDecrypter(e.block, iv)
	decrypter.XORKeyStream(encrypted, encrypted)
	return string(encrypted), nil
}
