package fake

import (
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/server"
)

// Settings are the extended settings for in-memory server.
type Settings struct {
	model.ServerSettings
}

// DefaultSettings are default server settings.
var DefaultSettings = Settings{
	ServerSettings: server.ServerSettings,
}

// NewServer creates new in-memory backend service.
func NewServer(settings Settings, options ...func(*server.Server) error) (model.Server, error) {
	dbComposer, err := NewComposer(settings)
	if err != nil {
		return nil, err
	}
	return server.NewServer(settings.ServerSettings, dbComposer, options...)
}
