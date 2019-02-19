package encryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"io/ioutil"

	"github.com/madappgang/identifo/model"
)

// EncryptitonService holds key.
type EncryptitonService struct {
	key []byte
}

// Encrypt encrypts plaintext with symmetric encryption method
func (es *EncryptitonService) Encrypt(plaintext []byte) ([]byte, error) {
	c, err := aes.NewCipher(es.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// Decrypt decrypts chiphertext with symmetric encryption method
func (es *EncryptitonService) Decrypt(ciphertext []byte) ([]byte, error) {
	c, err := aes.NewCipher(es.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// NewEncryptor creates a new encryptor
func NewEncryptor(keyPath string) (model.Encryptor, error) {
	es := EncryptitonService{}
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	es.key = key

	return es, nil
}
