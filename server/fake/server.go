package fake

import (
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/server"
)

//Settings is extended settings for memory server
type Settings struct {
	model.ServerSettings
}

//DefaultSettings default server settings
var DefaultSettings = Settings{
	ServerSettings: server.DefaultSettings,
}

//NewServer create memory backend service
func NewServer(setting Settings, options ...func(*server.Server) error) (model.Server, error) {
	dbComposer, err := NewComposer(setting)
	if err != nil {
		return nil, err
	}
	return server.NewServer(setting.ServerSettings, dbComposer, options...)
}
