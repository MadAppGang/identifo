package config

import (
	"fmt"
	"log"

	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage/fs"
	"github.com/madappgang/identifo/v2/storage/s3"
)

// InitConfigurationStorage initializes configuration storage.
func InitConfigurationStorage(config model.ConfigStorageSettings) (model.ConfigurationStorage, error) {
	switch config.Type {
	// case model.ConfigStorageTypeEtcd:
	// 	return etcd.NewConfigurationStorage(config)
	case model.ConfigStorageTypeS3:
		return s3.NewConfigurationStorage(config)
	case model.ConfigStorageTypeFile:
		return fs.NewConfigurationStorage(config)
	default:
		return nil, fmt.Errorf("config type is not supported")
	}
}

// DefaultStorage trying to create a default storage with default file
func DefaultStorage() (model.ConfigurationStorage, error) {
	return fs.NewDefaultConfigurationStorage()
}

func InitConfigurationStorageFromFlag(configFlag string) (model.ConfigurationStorage, error) {
	// ignore error to fall back to default if needed
	settings, settingsErr := model.ConfigStorageSettingsFromString(configFlag)
	configStorage, err := InitConfigurationStorage(settings)
	if err != nil || settingsErr != nil || configFlag == "" {
		log.Printf("Unable to init config using\n\tconfig string: %s\n\twith error: %v\nT",
			configFlag,
			err,
		)
		// Trying to fall back to default settings:
		log.Printf("Trying to load default settings from env variable 'SERVER_CONFIG_PATH' or default pathes.\n")
		configStorage, err = DefaultStorage()
		if err != nil {
			return nil, fmt.Errorf("Unable to load default config with error: %v", err)
		}
	}
	return configStorage, nil
}

func NewServerFromFlag(configFlag string, restartChan chan<- bool) (model.Server, error) {
	configStorage, err := InitConfigurationStorageFromFlag(configFlag)
	if err != nil {
		return nil, fmt.Errorf("Unable to load settings on start with error: %v ", err)
	}

	srv, err := NewServer(configStorage, restartChan)
	if err != nil {
		return nil, fmt.Errorf("Unable to create server with error: %v ", err)
	}

	return srv, nil
}
