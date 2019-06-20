package model

// ConfigurationWatcher is a global server configuration watcher.
type ConfigurationWatcher interface {
	Watch()
	WatchChan() chan interface{}
	Stop()
}
