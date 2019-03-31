package fake

import (
	"github.com/madappgang/identifo/mem"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/server"
)

// StorageProvider holds together all storage creators functions.
type StorageProvider struct {
	NewAppStorage   func() model.AppStorage
	NewUserStorage  func() model.UserStorage
	NewTokenStorage func() model.TokenStorage
}

var defaultSP = StorageProvider{
	NewAppStorage:   mem.NewAppStorage,
	NewUserStorage:  mem.NewUserStorage,
	NewTokenStorage: mem.NewTokenStorage,
}

// Settings are the extended settings for in-memory server.
type Settings struct {
	model.ServerSettings
	StorageProvider
}

// DefaultSettings are default server settings.
var DefaultSettings = Settings{
	ServerSettings:  server.DefaultSettings,
	StorageProvider: defaultSP,
}

// NewServer creates new in-memory backend service.
func NewServer(setting Settings, options ...func(*server.Server) error) (model.Server, error) {
	dbComposer, err := NewComposer(setting)
	if err != nil {
		return nil, err
	}
	return server.NewServer(setting.ServerSettings, dbComposer, options...)
}
