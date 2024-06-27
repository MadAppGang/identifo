package storage

import (
	"fmt"
	"log/slog"

	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage/fs"
	"github.com/madappgang/identifo/v2/storage/s3"
)

func NewKeyStorage(
	logger *slog.Logger,
	settings model.FileStorageSettings,
) (model.KeyStorage, error) {
	switch settings.Type {
	case model.FileStorageTypeLocal:
		return fs.NewKeyStorage(settings.Local)
	case model.FileStorageTypeS3:
		return s3.NewKeyStorage(logger, settings.S3)
	default:
		return nil, fmt.Errorf("unknown key storage type: %s", settings.Type)
	}
}
