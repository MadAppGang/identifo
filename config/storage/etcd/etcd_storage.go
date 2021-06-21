package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	keyStorageLocal "github.com/madappgang/identifo/config/key_storage/local"
	keyStorageS3 "github.com/madappgang/identifo/config/key_storage/s3"
	"github.com/madappgang/identifo/model"
	"go.etcd.io/etcd/clientv3"
)

const (
	timeoutPerRequest = 5 * time.Second
)

// ConfigurationStorage is an etcd-backed storage for server configuration.
type ConfigurationStorage struct {
	Client      *clientv3.Client
	settingsKey string
	keyStorage  model.KeyStorage
	cache       model.ServerSettings
	cached      bool
}

// NewConfigurationStorage creates new etcd-backed server config storage.
func NewConfigurationStorage(config model.ConfigStorageSettings) (*ConfigurationStorage, error) {
	log.Println("Loading server configuration from the etcd...")
	cfg := clientv3.Config{
		DialTimeout: timeoutPerRequest,
		Username:    config.Etcd.Username,
		Password:    config.Etcd.Password,
		Endpoints:   config.Etcd.Endpoints,
	}
	etcdClient, err := clientv3.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("Cannot not connect to etcd config storage: %s", err)
	}

	cs := ConfigurationStorage{
		Client:      etcdClient,
		settingsKey: config.Etcd.Key,
	}

	settings, err := cs.LoadServerSettings(true)
	if err != nil {
		return nil, fmt.Errorf("Cannot not load settings from etcd config storage: %s", err)
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

	cs.keyStorage = keyStorage
	cs.cached = false
	return &cs, nil
}

// WriteConfig write new configuration.
func (cs *ConfigurationStorage) WriteConfig(settings model.ServerSettings) error {
	// TODO: implement etcd update
	return fmt.Errorf("not supported")
	// _, err = cs.Client.Put(context.Background(), settings, strVal)
	// return err
}

// LoadServerSettings loads server configuration from configuration storage.
func (cs *ConfigurationStorage) LoadServerSettings(forceReload bool) (model.ServerSettings, error) {
	if !forceReload && cs.cached {
		return cs.cache, nil
	}
	res, err := cs.Client.Get(context.Background(), cs.settingsKey)
	if err != nil {
		return model.ServerSettings{},
			fmt.Errorf("Cannot get value by key %s: %s", cs.settingsKey, err)
	}
	if len(res.Kvs) == 0 {
		return model.ServerSettings{},
			fmt.Errorf("Etcd: No value for key %s", cs.settingsKey)
	}

	var settings model.ServerSettings
	err = json.Unmarshal(res.Kvs[0].Value, &settings)
	return settings, err
}

// InsertKeys inserts new public and private keys into the key storage.
func (cs *ConfigurationStorage) InsertKeys(keys *model.JWTKeys) error {
	if err := cs.keyStorage.InsertKeys(keys); err != nil {
		return err
	}
	return nil
}

// LoadKeys loads public and private keys from the key storage.
func (cs *ConfigurationStorage) LoadKeys(alg model.TokenSignatureAlgorithm) (*model.JWTKeys, error) {
	return cs.keyStorage.LoadKeys(alg)
}

// GetUpdateChan implements ConfigurationStorage interface.
func (cs *ConfigurationStorage) GetUpdateChan() chan interface{} {
	return make(chan interface{}, 1)
}

// CloseUpdateChan implements ConfigurationStorage interface.
func (cs *ConfigurationStorage) CloseUpdateChan() {}
