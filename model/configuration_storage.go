package model

// ConfigurationStorage stores server configuration.
type ConfigurationStorage interface {
	Insert(key string, value interface{}) error
	LoadServerSettings(*ServerSettings) error
}
