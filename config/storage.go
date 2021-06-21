package config

import (
	"fmt"

	"github.com/madappgang/identifo/config/storage/file"
	"github.com/madappgang/identifo/config/storage/s3"
	"github.com/madappgang/identifo/model"
)

// InitConfigurationStorage initializes configuration storage.
func InitConfigurationStorage(config model.ConfigStorageSettings) (model.ConfigurationStorage, error) {
	switch config.Type {
	// case model.ConfigStorageTypeEtcd:
	// 	return etcd.NewConfigurationStorage(config)
	case model.ConfigStorageTypeS3:
		return s3.NewConfigurationStorage(config)
	case model.ConfigStorageTypeFile:
		return file.NewConfigurationStorage(config)
	default:
		return nil, fmt.Errorf("config type is not supported")
	}
}

// DefaultStorage trying to create a default storage with default file
func DefaultStorage() (model.ConfigurationStorage, error) {
	return file.NewDefaultConfigurationStorage()
}
