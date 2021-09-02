package s3

import (
	"github.com/madappgang/identifo/model"
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
	_, err = ks.LoadKeys(model.TokenSignatureAlgorithmAuto)
	return err
}
