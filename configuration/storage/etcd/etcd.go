package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/madappgang/identifo/model"
	"go.etcd.io/etcd/clientv3"
)

const (
	etcdConnectionString = "http://127.0.0.1:2379"
	serverSettingsKey    = "identifo/server-settings"
	timeoutPerRequest    = 5 * time.Second
)

// ConfigurationStorage is an etcd-backed storage for server configuration.
type ConfigurationStorage struct {
	Client *clientv3.Client
}

// NewConfigurationStorage creates new etcd-backed server config storage.
func NewConfigurationStorage() (*ConfigurationStorage, error) {
	c, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdConnectionString},
		DialTimeout: timeoutPerRequest,
	})
	if err != nil {
		return nil, err
	}

	return &ConfigurationStorage{Client: c}, nil
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
			return fmt.Errorf("Cannot serialize struct %+v to string: %s", value, err)
		}
		strVal = string(out)
	}

	_, err = cs.Client.Put(context.Background(), key, strVal)
	return err
}

// LoadServerSettings loads server configuration from configuration storage.
func (cs *ConfigurationStorage) LoadServerSettings(settings *model.ServerSettings) error {
	res, err := cs.Client.Get(context.Background(), serverSettingsKey)
	if err != nil {
		return fmt.Errorf("Cannot get value by key %s: %s", serverSettingsKey, err)
	}

	err = json.Unmarshal(res.Kvs[0].Value, settings)
	return err
}
