package fs

import (
	"fmt"
	"io"
	"os"

	"github.com/madappgang/identifo/jwt"
	"github.com/madappgang/identifo/model"
)

var supportedSignatureAlgorithms = [2]model.TokenSignatureAlgorithm{
	model.TokenSignatureAlgorithmES256,
	model.TokenSignatureAlgorithmRS256,
}

// KeyStorage is a wrapper over public and private key files.
type KeyStorage struct {
	PublicKeyPath  string
	PrivateKeyPath string
}

// NewKeyStorage creates and returns new key files storage.
func NewKeyStorage(settings model.KeyStorageFileSettings) (*KeyStorage, error) {
	return &KeyStorage{
		PrivateKeyPath: settings.PrivateKeyPath,
		PublicKeyPath:  settings.PublicKeyPath,
	}, nil
}

// InsertKeys inserts public and private keys.
func (ks *KeyStorage) ReplaceKeys(keys model.JWTKeys) error {
	if keys.Private == nil || keys.Public == nil {
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
func (ks *KeyStorage) LoadKeys(alg model.TokenSignatureAlgorithm) (model.JWTKeys, error) {
	if _, err := os.Stat(ks.PublicKeyPath); err != nil {
		if os.IsNotExist(err) {
			wd, wdErr := os.Getwd()
			if wdErr != nil {
				return model.JWTKeys{}, fmt.Errorf("Public key not found. Also, cannot get working directory: %s", wdErr)
			}
			return model.JWTKeys{}, fmt.Errorf("Public key file not found. Working directory: %s, key path: %s", wd, ks.PublicKeyPath)
		}
		return model.JWTKeys{}, fmt.Errorf("Error while checking public key existence. %s", err)
	}

	if _, err := os.Stat(ks.PrivateKeyPath); err != nil {
		if os.IsNotExist(err) {
			wd, wdErr := os.Getwd()
			if wdErr != nil {
				return model.JWTKeys{}, fmt.Errorf("Private key not found. Also, cannot get working directory: %s", wdErr)
			}
			return model.JWTKeys{}, fmt.Errorf("Private key file not found. Working directory: %s, key path: %s", wd, ks.PrivateKeyPath)
		}
		return model.JWTKeys{}, fmt.Errorf("Error while checking private key existence. %s", err)
	}

	keys := model.JWTKeys{}
	if alg != model.TokenSignatureAlgorithmAuto {
		return ks.loadKeys(alg)
	}

	// Trying to guess algorithm
	var err error
	for _, a := range supportedSignatureAlgorithms {
		keys, err = ks.loadKeys(a)
		if err == nil {
			break
		}
	}
	return keys, err
}

func (ks *KeyStorage) loadKeys(alg model.TokenSignatureAlgorithm) (model.JWTKeys, error) {
	keys := model.JWTKeys{}

	privateKey, err := jwt.LoadPrivateKeyFromPEM(ks.PrivateKeyPath, alg)
	if err != nil {
		return keys, fmt.Errorf("Cannot load private key: %s", err)
	}

	publicKey, err := jwt.LoadPublicKeyFromPEM(ks.PublicKeyPath, alg)
	if err != nil {
		return keys, fmt.Errorf("Cannot load public key: %s", err)
	}
	keys.Private = privateKey
	keys.Public = publicKey
	keys.Algorithm = alg
	return keys, nil
}
