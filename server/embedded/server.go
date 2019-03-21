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

// ServerSettings returns  server settings.
var ServerSettings = Settings{
	ServerSettings: server.ServerSettings,
	DBPath:         "db.db",
}

// NewServer creates new backend service with BoltDB support.
func NewServer(settings Settings, options ...func(*server.Server) error) (model.Server, error) {
	dbComposer, err := NewComposer(settings)
	if err != nil {
		return nil, err
	}
	return server.NewServer(settings.ServerSettings, dbComposer, options...)
}
