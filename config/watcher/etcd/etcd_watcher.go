package etcd

import (
	"context"

	"github.com/madappgang/identifo/config/storage/etcd"
	"go.etcd.io/etcd/clientv3"
)

// ConfigurationWatcher wraps etcd client.
type ConfigurationWatcher struct {
	Client            *clientv3.Client
	watchChan         chan interface{}
	serverSettingsKey string
}

// NewConfigurationWatcher creates and returns new etcd-backed configuration watcher.
func NewConfigurationWatcher(configStorage *etcd.ConfigurationStorage, settingsKey string, watchChan chan interface{}) (*ConfigurationWatcher, error) {
	return &ConfigurationWatcher{
		Client:            configStorage.Client,
		watchChan:         watchChan,
		serverSettingsKey: settingsKey,
	}, nil
}

// Watch watches for configuration updates.
func (cw *ConfigurationWatcher) Watch() {
	internalWatchChan := cw.Client.Watch(context.Background(), cw.serverSettingsKey)

	go func() {
		for event := range internalWatchChan {
			if event.Canceled {
				return
			}
			cw.watchChan <- event
		}
	}()
}

// WatchChan returns watcher's event channel.
func (cw *ConfigurationWatcher) WatchChan() chan interface{} {
	return cw.watchChan
}

// Stop stops watcher.
func (cw *ConfigurationWatcher) Stop() {
	cw.Client.Close()
}
