package model

// ConfigurationStorage stores server configuration.
type ConfigurationStorage interface {
	WriteConfig(ServerSettings) error
	LoadServerSettings(forceReload bool) (ServerSettings, error)
	GetUpdateChan() chan interface{}
	CloseUpdateChan()
}
