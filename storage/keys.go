package storage

import (
	"fmt"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/storage/fs"
	"github.com/madappgang/identifo/storage/s3"
)

func NewKeyStorage(settings model.KeyStorageSettings) (model.KeyStorage, error) {
	switch settings.Type {
	case model.KeyStorageTypeLocal:
		return fs.NewKeyStorage(settings)
	case model.KeyStorageTypeS3:
		return s3.NewKeyStorage(settings)
	default:
		return nil, fmt.Errorf("unknown key storage type: %s", settings.Type)
	}
}
