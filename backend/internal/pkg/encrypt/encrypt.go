package encrypt

type Encrypt interface {
	Encrypt(plaintext []byte) (string, error)
	Decrypt(ciphertext string) (string, error)
}
