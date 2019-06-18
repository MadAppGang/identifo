package mock

// ConfigurationWatcher is a mock for real config watcher.
type ConfigurationWatcher struct{}

// NewConfigurationWatcher creates and returns new mocked configuration watcher.
func NewConfigurationWatcher() (*ConfigurationWatcher, error) {
	return &ConfigurationWatcher{}, nil
}

// Watch does basically nothing.
func (cw *ConfigurationWatcher) Watch() {}

// WatchChan returns non-nil yet useless channel of empty interfaces.
func (cw *ConfigurationWatcher) WatchChan() chan interface{} {
	return make(chan interface{}, 1)
}

// Stop does basically nothing.
func (cw *ConfigurationWatcher) Stop() {}
