package dynamodb

import (
	"github.com/madappgang/identifo/dynamodb"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/server"
)

// StorageProvider holds together all storage creators functions.
type StorageProvider struct {
	NewAppStorage   func(*dynamodb.DB) (model.AppStorage, error)
	NewUserStorage  func(*dynamodb.DB) (model.UserStorage, error)
	NewTokenStorage func(*dynamodb.DB) (model.TokenStorage, error)
}

var defaultSP = StorageProvider{
	NewAppStorage:   dynamodb.NewAppStorage,
	NewUserStorage:  dynamodb.NewUserStorage,
	NewTokenStorage: dynamodb.NewTokenStorage,
}

// Settings are the extended settings for DynamoDB server.
type Settings struct {
	model.ServerSettings
	DBEndpoint string
	DBRegion   string
	StorageProvider
}

// DefaultSettings are default server settings.
var DefaultSettings = Settings{
	ServerSettings:  server.DefaultSettings,
	DBEndpoint:      "",
	DBRegion:        "",
	StorageProvider: defaultSP,
}

// NewServer creates new backend service with DynamoDB support.
func NewServer(setting Settings, options ...func(*server.Server) error) (model.Server, error) {
	dbComposer, err := NewComposer(setting)
	if err != nil {
		return nil, err
	}
	return server.NewServer(setting.ServerSettings, dbComposer, options...)
}
