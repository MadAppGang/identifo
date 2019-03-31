package mgo

import (
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/mongo"
	"github.com/madappgang/identifo/server"
)

// StorageProvider holds together all storage creators functions.
type StorageProvider struct {
	NewAppStorage   func(*mongo.DB) (model.AppStorage, error)
	NewUserStorage  func(*mongo.DB) (model.UserStorage, error)
	NewTokenStorage func(*mongo.DB) (model.TokenStorage, error)
}

var defaultSP = StorageProvider{
	NewAppStorage:   mongo.NewAppStorage,
	NewUserStorage:  mongo.NewUserStorage,
	NewTokenStorage: mongo.NewTokenStorage,
}

// Settings are the extended settings for MongoDB server.
type Settings struct {
	model.ServerSettings
	DBEndpoint string
	DBName     string
	StorageProvider
}

// DefaultSettings are default server settings.
var DefaultSettings = Settings{
	ServerSettings:  server.DefaultSettings,
	DBEndpoint:      "localhost:27017",
	DBName:          "identifo",
	StorageProvider: defaultSP,
}

// NewServer creates new backend service with MongoDB support.
func NewServer(setting Settings, options ...func(*server.Server) error) (model.Server, error) {
	dbComposer, err := NewComposer(setting)
	if err != nil {
		return nil, err
	}
	return server.NewServer(setting.ServerSettings, dbComposer, options...)
}
