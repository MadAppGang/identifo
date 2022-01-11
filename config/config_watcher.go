package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/madappgang/identifo/v2/config/storage/file"
	"github.com/madappgang/identifo/v2/config/storage/s3"
	"github.com/madappgang/identifo/v2/model"
)

const defaultS3PollInterval = time.Minute // poll new updates every minute

func NewConfigWatcher(settings model.ConfigStorageSettings) (model.ConfigurationWatcher, error) {
	switch settings.Type {
	case model.ConfigStorageTypeS3:
		if settings.S3 == nil {
			return nil, errors.New("empty storage settings for S3 type")
		}
		return s3.NewPollWatcher(*settings.S3, defaultS3PollInterval)
	case model.ConfigStorageTypeFile:
		if settings.File == nil {
			return nil, errors.New("empty storage settings for File storage type")
		}
		return file.NewWatcher(settings.File.FileName), nil
	}
	return nil, fmt.Errorf("Unsupported config storage type: %v", settings.Type)
}
