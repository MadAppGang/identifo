package mgo

import (
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/server"
)

// Settings are the extended settings for MongoDB server.
type Settings struct {
	model.ServerSettings
	DBEndpoint string
	DBName     string
}

// DefaultSettings are default server settings.
var DefaultSettings = Settings{
	ServerSettings: server.DefaultSettings,
	DBEndpoint:     "localhost:27017",
	DBName:         "identifo",
}

// NewServer creates new backend service with MongoDB support.
func NewServer(setting Settings, options ...func(*server.Server) error) (model.Server, error) {
	dbComposer, err := NewComposer(setting)
	if err != nil {
		return nil, err
	}
	return server.NewServer(setting.ServerSettings, dbComposer, options...)
}
