package fs

import (
	"github.com/madappgang/identifo/v2/model"
)

type KeyStorageConnectionTester struct {
	settings model.FileStorageLocal
}

// NewConnectionTester creates a BoltDB connection tester

func NewKeyStorageConnectionTester(settings model.FileStorageLocal) model.ConnectionTester {
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
