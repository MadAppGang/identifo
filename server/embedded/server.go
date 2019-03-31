package embedded

import (
	"github.com/boltdb/bolt"
	"github.com/madappgang/identifo/boltdb"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/server"
)

// StorageProvider holds together all storage creators functions.
type StorageProvider struct {
	NewAppStorage   func(*bolt.DB) (model.AppStorage, error)
	NewUserStorage  func(*bolt.DB) (model.UserStorage, error)
	NewTokenStorage func(*bolt.DB) (model.TokenStorage, error)
}

var defaultSP = StorageProvider{
	NewAppStorage:   boltdb.NewAppStorage,
	NewUserStorage:  boltdb.NewUserStorage,
	NewTokenStorage: boltdb.NewTokenStorage,
}

// Settings are the extended settings for BoltDB server.
type Settings struct {
	model.ServerSettings
	DBPath string
	StorageProvider
}

// DefaultSettings are default server settings.
var DefaultSettings = Settings{
	ServerSettings:  server.DefaultSettings,
	DBPath:          "db.db",
	StorageProvider: defaultSP,
}

// NewServer creates new backend service with BoltDB support.
func NewServer(setting Settings, options ...func(*server.Server) error) (model.Server, error) {
	dbComposer, err := NewComposer(setting)
	if err != nil {
		return nil, err
	}
	return server.NewServer(setting.ServerSettings, dbComposer, options...)
}
