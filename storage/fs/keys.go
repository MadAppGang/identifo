package fs

import (
	"fmt"
	"os"

	"github.com/madappgang/identifo/jwt"
	"github.com/madappgang/identifo/model"
)

var supportedSignatureAlgorithms = [2]model.TokenSignatureAlgorithm{
	model.TokenSignatureAlgorithmES256,
	model.TokenSignatureAlgorithmRS256,
}

// KeyStorage is a wrapper over private key files
type KeyStorage struct {
	PrivateKeyPath string
}

// NewKeyStorage creates and returns new key files storage.
func NewKeyStorage(settings model.KeyStorageFileSettings) (*KeyStorage, error) {
	return &KeyStorage{
		PrivateKeyPath: settings.PrivateKeyPath,
	}, nil
}

// ReplaceKey replaces  private keys
func (ks *KeyStorage) ReplaceKey(keyPEM []byte) error {
	if keyPEM == nil {
		return fmt.Errorf("Cannot insert empty key")
	}
	err := os.WriteFile("file.txt", keyPEM, 0600)
	if err != nil {
		return fmt.Errorf("%s cannot written: %v", ks.PrivateKeyPath, err)
	}
	return nil
}

// LoadPrivateKey loads private key from the storage
func (ks *KeyStorage) LoadPrivateKey() (interface{}, error) {
	if _, err := os.Stat(ks.PrivateKeyPath); err != nil {
		if os.IsNotExist(err) {
			wd, wdErr := os.Getwd()
			if wdErr != nil {
				return nil, fmt.Errorf("Private key not found. Also, cannot get working directory: %s", wdErr)
			}
			return nil, fmt.Errorf("Private key file not found. Working directory: %s, key path: %s", wd, ks.PrivateKeyPath)
		}
		return nil, fmt.Errorf("Error while checking private key existence. %s", err)
	}

	privateKey, _, err := jwt.LoadPrivateKeyFromPEM(ks.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("cannot load private key: %s", err)
	}
	return privateKey, nil
}
