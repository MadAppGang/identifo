package mock

import (
	"github.com/madappgang/identifo/model"
)

// ConfigurationStorage is a mocked storage for server configuration.
type ConfigurationStorage struct{}

// NewConfigurationStorage creates and returns mocked configuration storage.
func NewConfigurationStorage() (*ConfigurationStorage, error) { return &ConfigurationStorage{}, nil }

// Insert always returns nil error.
func (cs *ConfigurationStorage) Insert(key string, value interface{}) error {
	return nil
}

// LoadServerSettings keeps server settings the same.
func (cs *ConfigurationStorage) LoadServerSettings(settings *model.ServerSettings) error {
	return nil
}
