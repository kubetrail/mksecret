package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	minPassphraseLen = 8
)

// NewAesKeyFromPassphrase generates new AES key deterministically using input key
func NewAesKeyFromPassphrase(passphrase []byte) ([]byte, error) {
	if len(passphrase) < minPassphraseLen {
		return nil, fmt.Errorf("passphrase length needs to be at least 8")
	}

	salt := md5.Sum(passphrase)
	key := pbkdf2.Key(passphrase, salt[:], 4096, 32, sha256.New)
	return key, nil
}

// EncryptWithAesKey encrypts data using AES key
func EncryptWithAesKey(data, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		err := fmt.Errorf("could not create a new aes cipher: %w", err)
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		err := fmt.Errorf("could not create new gcm from cipher: %w", err)
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		err := fmt.Errorf("could not populate nonce: %w", err)
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

// DecryptWithAesKey decrypts data using AES key
func DecryptWithAesKey(data, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		err := fmt.Errorf("could not create a new aes cipher: %w", err)
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		err := fmt.Errorf("could not create new gcm from cipher: %w", err)
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		err := fmt.Errorf("invalid cipher text, length less than nonce")
		return nil, err
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		err := fmt.Errorf("could not decrypt cipher text: %w", err)
		return nil, err
	}

	return plaintext, nil
}
