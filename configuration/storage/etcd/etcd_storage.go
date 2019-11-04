package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"time"

	keyStorageLocal "github.com/madappgang/identifo/configuration/key_storage/local"
	keyStorageS3 "github.com/madappgang/identifo/configuration/key_storage/s3"
	configStorageFile "github.com/madappgang/identifo/configuration/storage/file"
	ijwt "github.com/madappgang/identifo/jwt"
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
	settingsKey string
	keyStorage  model.KeyStorage
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

	settingsKey := settings.SettingsKey
	// Init file storage for config replication.
	settings.SettingsKey = serverConfigPath
	fileStorage, err := configStorageFile.NewConfigurationStorage(settings)
	if err != nil {
		return nil, err
	}

	var keyStorage model.KeyStorage

	switch settings.KeyStorage.Type {
	case model.KeyStorageTypeLocal:
		keyStorage, err = keyStorageLocal.NewKeyStorage(settings.KeyStorage)
	case model.KeyStorageTypeS3:
		keyStorage, err = keyStorageS3.NewKeyStorage(settings.KeyStorage)
	default:
		return nil, fmt.Errorf("Unknown key storage type: %s", settings.KeyStorage.Type)
	}
	if err != nil {
		return nil, err
	}

	return &ConfigurationStorage{
		Client:      c,
		FileStorage: fileStorage,
		settingsKey: settingsKey,
		keyStorage:  keyStorage,
	}, nil
}

// InsertConfig inserts key-value pair to configuration storage.
func (cs *ConfigurationStorage) InsertConfig(key string, value interface{}) error {
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

	if key == "" && strVal == "" {
		go cs.idleInsertConfig()
		return nil
	}

	_, err = cs.Client.Put(context.Background(), key, strVal)
	if err == nil {
		// Also update file.
		go func() {
			if fileErr := cs.FileStorage.InsertConfig(cs.FileStorage.ServerConfigPath, value); fileErr != nil {
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
	res, err := cs.Client.Get(context.Background(), cs.settingsKey)
	if err != nil {
		return fmt.Errorf("Cannot get value by key %s: %s", cs.settingsKey, err)
	}
	if len(res.Kvs) == 0 {
		return fmt.Errorf("Etcd: No value for key %s", cs.settingsKey)
	}

	err = json.Unmarshal(res.Kvs[0].Value, settings)
	return err
}

// InsertKeys inserts new public and private keys into the key storage.
func (cs *ConfigurationStorage) InsertKeys(keys *model.JWTKeys) error {
	if err := cs.keyStorage.InsertKeys(keys); err != nil {
		return err
	}
	return nil
}

// LoadKeys loads public and private keys from the key storage.
func (cs *ConfigurationStorage) LoadKeys(alg ijwt.TokenSignatureAlgorithm) (*model.JWTKeys, error) {
	return cs.keyStorage.LoadKeys(alg)
}

// GetUpdateChan implements ConfigurationStorage interface.
func (cs *ConfigurationStorage) GetUpdateChan() chan interface{} {
	return make(chan interface{}, 1)
}

// CloseUpdateChan implements ConfigurationStorage interface.
func (cs *ConfigurationStorage) CloseUpdateChan() {}

// idleInsertConfig inserts existing settings.
func (cs *ConfigurationStorage) idleInsertConfig() {
	key := cs.settingsKey
	settings := new(model.ServerSettings)
	if err := cs.LoadServerSettings(settings); err != nil {
		log.Println("Error while idle config insert: could not load server settings.", err)
		return
	}
	if key == "" {
		log.Println("Error while idle config insert: empty key.")
		return
	}
	if err := cs.InsertConfig(key, settings); err != nil {
		log.Println("Error while idle config insert.", err)
		return
	}
}
