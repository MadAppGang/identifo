package generic

import (
	"github.com/madappgang/identifo/model"
)

// ConfigurationWatcher is a storage-agnostic config watcher.
type ConfigurationWatcher struct {
	Storage           model.ConfigurationStorage
	watchChan         chan interface{}
	serverSettingsKey string
}

// NewConfigurationWatcher creates and returns new storage-agnostic configuration watcher.
func NewConfigurationWatcher(configStorage model.ConfigurationStorage, settingsKey string, watchChan chan interface{}) (*ConfigurationWatcher, error) {
	return &ConfigurationWatcher{
		Storage:           configStorage,
		serverSettingsKey: settingsKey,
		watchChan:         watchChan,
	}, nil
}

// Watch watches for configuration updates.
func (cw *ConfigurationWatcher) Watch() {
	internalWatchChan := cw.Storage.GetUpdateChan()
	go func() {
		for event := range internalWatchChan {
			cw.watchChan <- event
		}
	}()
}

// WatchChan returns watcher's event channel.
func (cw *ConfigurationWatcher) WatchChan() chan interface{} {
	return cw.watchChan
}

// Stop stops listening on config updates.
func (cw *ConfigurationWatcher) Stop() {
	cw.Storage.CloseUpdateChan()
	close(cw.WatchChan())
}
