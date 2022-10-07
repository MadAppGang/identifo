package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage"
)

const defaultS3PollInterval = time.Minute // poll new updates every minute

func NewConfigWatcher(settings model.FileStorageSettings) (model.ConfigurationWatcher, error) {
	filename := settings.FileName()
	switch settings.Type {
	case model.FileStorageTypeS3:
		if len(settings.S3.Bucket) == 0 {
			return nil, errors.New("empty storage settings for S3 type")
		}
		settings.S3.Key = settings.Dir() // remove filename from key to keep dir only
	case model.FileStorageTypeLocal:
		if len(settings.Local.Path) == 0 {
			return nil, errors.New("empty storage settings for File storage type")
		}
		settings.Local.Path = settings.Dir() // remove filename from key to keep dir only
	default:
		return nil, fmt.Errorf("Unsupported config storage type: %v", settings.Type)
	}

	fs, err := storage.NewFS(settings)
	if err != nil {
		return nil, err
	}
	watcher := storage.NewFSWatcher(fs, []string{filename}, defaultS3PollInterval)
	return watcher, nil
}
