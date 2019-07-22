package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	configStorageFile "github.com/madappgang/identifo/configuration/storage/file"
	"github.com/madappgang/identifo/model"
	"go.etcd.io/etcd/clientv3"
)

const (
	defaultEtcdConnectionString = "http://127.0.0.1:2379"
	timeoutPerRequest           = 5 * time.Second
)

// ConfigurationStorage is an etcd-backed storage for server configuration.
type ConfigurationStorage struct {
	Client      *clientv3.Client
	FileStorage *configStorageFile.ConfigurationStorage
}

// NewConfigurationStorage creates new etcd-backed server config storage.
func NewConfigurationStorage(settings model.ConfigurationStorageSettings, serverConfigPath string) (*ConfigurationStorage, error) {
	if settings.SettingsKey == "" {
		return nil, fmt.Errorf("Empty server settings key for etcd")
	}

	if settings.Endpoints == nil {
		settings.Endpoints = []string{defaultEtcdConnectionString}
	}

	c, err := clientv3.New(clientv3.Config{
		Endpoints:   settings.Endpoints,
		DialTimeout: timeoutPerRequest,
	})

	if err != nil {
		return nil, err
	}

	// Init file storage for config replication.
	settings.SettingsKey = serverConfigPath
	fileStorage, err := configStorageFile.NewConfigurationStorage(settings)
	if err != nil {
		return nil, err
	}

	return &ConfigurationStorage{
		Client:      c,
		FileStorage: fileStorage,
	}, nil
}

// Insert inserts key-value pair to configuration storage.
func (cs *ConfigurationStorage) Insert(key string, value interface{}) error {
	var strVal string
	var err error

	switch reflect.TypeOf(value).Kind() {
	case reflect.String:
		strVal = value.(string)
	case reflect.Ptr:
		out, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("Cannot serialize pointer %v to string: %s", value, err)
		}
		strVal = string(out)
	}

	_, err = cs.Client.Put(context.Background(), key, strVal)
	if err == nil {
		// Also update file.
		go func() {
			if fileErr := cs.FileStorage.Insert(cs.FileStorage.ServerConfigPath, value); fileErr != nil {
				fmt.Println("Could not replicate settings in file: ", fileErr)
			} else {
				fmt.Println("Successfully replicated settings in file")
			}
		}()
	}
	return err
}

// LoadServerSettings loads server configuration from configuration storage.
func (cs *ConfigurationStorage) LoadServerSettings(settings *model.ServerSettings) error {
	res, err := cs.Client.Get(context.Background(), settings.ConfigurationStorage.SettingsKey)
	if err != nil {
		return fmt.Errorf("Cannot get value by key %s: %s", settings.ConfigurationStorage.SettingsKey, err)
	}

	err = json.Unmarshal(res.Kvs[0].Value, settings)
	return err
}

// GetUpdateChan implements ConfigurationStorage interface.
func (cs *ConfigurationStorage) GetUpdateChan() chan interface{} {
	return make(chan interface{}, 1)
}

// CloseUpdateChan implements ConfigurationStorage interface.
func (cs *ConfigurationStorage) CloseUpdateChan() {}
