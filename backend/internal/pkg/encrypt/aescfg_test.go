package encrypt

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const testEncryptionKey = "development-only-change-me-key!!"

func TestAESAuthenticatedRoundTrip(t *testing.T) {
	encrypter, err := NewAesCfb(testEncryptionKey)
	require.NoError(t, err)

	cipherText, err := encrypter.Encrypt([]byte("authenticated payload"))
	require.NoError(t, err)
	require.True(t, strings.HasPrefix(cipherText, authenticatedCipherPrefix))

	plainText, err := encrypter.Decrypt(cipherText)
	require.NoError(t, err)
	require.Equal(t, "authenticated payload", plainText)
}

func TestAESRejectsTamperedCiphertext(t *testing.T) {
	encrypter, err := NewAesCfb(testEncryptionKey)
	require.NoError(t, err)
	cipherText, err := encrypter.Encrypt([]byte("authenticated payload"))
	require.NoError(t, err)

	last := cipherText[len(cipherText)-1]
	replacement := byte('A')
	if last == replacement {
		replacement = 'B'
	}
	_, err = encrypter.Decrypt(cipherText[:len(cipherText)-1] + string(replacement))
	require.Error(t, err)
}

func TestAESDecryptsLegacyCFBPayload(t *testing.T) {
	encrypter, err := NewAesCfb(testEncryptionKey)
	require.NoError(t, err)

	plainText, err := encrypter.Decrypt("AAAAAAAAAAAAAAAAAAAAAMUoqIsDvy/8pgYbNX6k")
	require.NoError(t, err)
	require.Equal(t, "legacy payload", plainText)
}
