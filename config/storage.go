package config

import (
	"fmt"
	"log"
	"os"

	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage/fs"
	"github.com/madappgang/identifo/v2/storage/s3"
)

// InitConfigurationStorage initializes configuration storage.
func InitConfigurationStorage(config model.FileStorageSettings) (model.ConfigurationStorage, error) {
	switch config.Type {
	// case model.ConfigStorageTypeEtcd:
	// 	return etcd.NewConfigurationStorage(config)
	case model.FileStorageTypeS3:
		return s3.NewConfigurationStorage(config)
	case model.FileStorageTypeLocal:
		return fs.NewConfigurationStorage(config)
	default:
		return nil, fmt.Errorf("config type is not supported")
	}
}

// DefaultStorage trying to create a default storage with default file
func DefaultStorage() model.ConfigurationStorage {
	return fs.NewDefaultConfigurationStorage()
}

func InitConfigurationStorageFromFlag(configFlag string) (model.ConfigurationStorage, error) {
	// trying to get server settings from env variable
	if len(configFlag) == 0 {
		configFlag = os.Getenv(model.IdentifoConfigPathEnvName)
	}
	// if we have no config flag available and not env variable set, just load default config file
	if len(configFlag) == 0 {
		log.Println("Config Storage: not config flag specified, I am loading default build in config file")
		return DefaultStorage(), nil
	}

	// if config settings are invalid and not empty we should stop the app
	// as it means the service is misconfigured and could not works at all
	settings, err := model.ConfigStorageSettingsFromString(configFlag)
	if err != nil {
		return nil, fmt.Errorf("Unable to init config using\n\tconfig string: %s\n\twith error: %v\nT",
			configFlag,
			err,
		)
	}

	configStorage, err := InitConfigurationStorage(settings)
	if err != nil {
		return nil, fmt.Errorf("Unable to init config using\n\tconfig string: %s\n\twith error: %v\nT",
			configFlag,
			err)
	}
	return configStorage, nil
}

func NewServerFromFlag(configFlag string, restartChan chan<- bool) (model.Server, error) {
	configStorage, err := InitConfigurationStorageFromFlag(configFlag)
	if err != nil {
		return nil, fmt.Errorf("Unable to load settings on start with error: %v ", err)
	}

	// config storage should:
	// load default settings if the originals settings file is unavailable
	// continue check desired config location in case the proper config file appear
	// validate config settings
	// if server settings are invalid - load fallback settings and mark itself ad invalid
	// continue to listen for file location and reload it in case it changed
	//
	// this means that the location for settings is valid
	// but settings are invalid or unreachable
	// both things could be fixed while the app is running
	// that is why Identifo should run, letting the admin to fix it
	srv, err := NewServer(configStorage, restartChan)
	if err != nil {
		return nil, fmt.Errorf("Unable to create server with error: %v ", err)
	}

	return srv, nil
}
