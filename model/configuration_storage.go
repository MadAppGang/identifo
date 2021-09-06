package model

// ConfigurationStorage stores server configuration.
type ConfigurationStorage interface {
	WriteConfig(ServerSettings) error
	LoadServerSettings(forceReload bool) (ServerSettings, error)
	GetUpdateChan() chan interface{}
	CloseUpdateChan()

	// ForceReloadOnWriteConfig function returns the bool
	// if true - after WriteConfig we need to force reload server to apply the changes
	// if false - we don't need force reload server, because the watcher will reload the server instantly
	// for example S3 storage uses 1  mins polling, and to apply  new changes instantly we need to force restart the server
	// for file storage we don't need to force reload it. The file watcher will notify about file change instantly
	ForceReloadOnWriteConfig() bool
}
