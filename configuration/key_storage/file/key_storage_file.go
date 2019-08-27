package file

import (
	"fmt"
	"io"
	"os"

	ijwt "github.com/madappgang/identifo/jwt"
	"github.com/madappgang/identifo/model"
)

var supportedSignatureAlgorithms = [2]ijwt.TokenSignatureAlgorithm{ijwt.TokenSignatureAlgorithmES256, ijwt.TokenSignatureAlgorithmRS256}

// KeyStorage is a wrapper over public and private key files.
type KeyStorage struct {
	PublicKeyPath  string
	PrivateKeyPath string
}

// NewKeyStorage creates and returns new key files storage.
func NewKeyStorage(settings model.KeyStorageSettings) (*KeyStorage, error) {
	return &KeyStorage{
		PublicKeyPath:  settings.PublicKey,
		PrivateKeyPath: settings.PrivateKey,
	}, nil
}

// InsertKeys inserts public and private keys.
func (ks *KeyStorage) InsertKeys(keys *model.JWTKeys) error {
	if keys == nil || keys.Private == nil || keys.Public == nil {
		return fmt.Errorf("Cannot insert empty key(s)")
	}
	keysMap := map[string]interface{}{
		ks.PrivateKeyPath: keys.Private,
		ks.PublicKeyPath:  keys.Public,
	}

	for name, file := range keysMap {
		reader, ok := file.(io.Reader)
		if !ok {
			return fmt.Errorf("%s cannot be read", name)
		}

		keyFile, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return fmt.Errorf("Cannot open key file: %s", err.Error())
		}
		defer keyFile.Close()

		if _, err := io.Copy(keyFile, reader); err != nil {
			return fmt.Errorf("Cannot copy key file contents: %s", err.Error())
		}
	}
	return nil
}

// LoadKeys loads keys from the key storage.
func (ks *KeyStorage) LoadKeys(alg ijwt.TokenSignatureAlgorithm) (*model.JWTKeys, error) {
	if _, err := os.Stat(ks.PublicKeyPath); err != nil {
		return nil, fmt.Errorf("Public key file not found")
	}
	if _, err := os.Stat(ks.PrivateKeyPath); err != nil {
		return nil, fmt.Errorf("Private key file not found")
	}

	keys := new(model.JWTKeys)

	if alg != ijwt.TokenSignatureAlgorithmAuto {
		if err := ks.loadKeys(alg, keys); err != nil {
			return nil, err
		}
		return keys, nil
	}

	// Trying to guess algo.
	var err error
	for _, a := range supportedSignatureAlgorithms {
		if err = ks.loadKeys(a, keys); err == nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (ks *KeyStorage) loadKeys(alg ijwt.TokenSignatureAlgorithm, keys *model.JWTKeys) error {
	privateKey, err := ijwt.LoadPrivateKeyFromPEM(ks.PrivateKeyPath, alg)
	if err != nil {
		return fmt.Errorf("Cannot load private key: %s", err)
	}
	publicKey, err := ijwt.LoadPublicKeyFromPEM(ks.PublicKeyPath, alg)
	if err != nil {
		return fmt.Errorf("Cannot load public key: %s", err)
	}
	keys.Private = privateKey
	keys.Public = publicKey
	keys.Algorithm = alg
	return nil
}
