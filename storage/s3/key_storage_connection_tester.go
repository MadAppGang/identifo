package s3

import (
	"github.com/madappgang/identifo/v2/model"
)

type KeyStorageConnectionTester struct {
	settings model.S3KeyStorageSettings
}

// NewConnectionTester creates a BoltDB connection tester

func NewKeyStorageConnectionTester(settings model.S3KeyStorageSettings) model.ConnectionTester {
	return &KeyStorageConnectionTester{settings: settings}
}

func (ct *KeyStorageConnectionTester) Connect() error {
	// let's try to load keys from the storage, if we can - it means
	ks, err := NewKeyStorage(ct.settings)
	if err != nil {
		return err
	}
	_, err = ks.LoadPrivateKey()
	return err
}
