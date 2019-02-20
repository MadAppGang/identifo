package embedded

import (
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/server"
)

// Settings are the extended settings for BoltDB server.
type Settings struct {
	model.ServerSettings
	DBPath string
}

// DefaultSettings are default server settings.
var DefaultSettings = Settings{
	ServerSettings: server.DefaultSettings,
	DBPath:         "db.db",
}

// NewServer creates new backend service with BoltDB support.
func NewServer(setting Settings, options ...func(*server.Server) error) (model.Server, error) {
	dbComposer, err := NewComposer(setting)
	if err != nil {
		return nil, err
	}
	return server.NewServer(setting.ServerSettings, dbComposer, options...)
}
