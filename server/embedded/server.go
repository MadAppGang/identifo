package embedded

import (
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/server"
)

//Settings is extended settings for BoltDB serve
type Settings struct {
	model.ServerSettings
	DBPath string
}

//DefaultSettings default server settings
var DefaultSettings = Settings{
	ServerSettings: server.DefaultSettings,
	DBPath:         "db.db",
}

//NewServer creates BoltDB backend service
func NewServer(setting Settings, options ...func(*server.Server) error) (model.Server, error) {
	dbComposer, err := NewComposer(setting)
	if err != nil {
		return nil, err
	}
	return server.NewServer(setting.ServerSettings, dbComposer, options...)
}
