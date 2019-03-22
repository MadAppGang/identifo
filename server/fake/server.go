package fake

import (
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/server"
)

// DefaultSettings are default server settings.
var DefaultSettings = server.ServerSettings

// NewServer creates new in-memory backend service.
func NewServer(settings model.ServerSettings, options ...func(*server.Server) error) (model.Server, error) {
	dbComposer, err := NewComposer(settings)
	if err != nil {
		return nil, err
	}
	return server.NewServer(settings, dbComposer, options...)
}
