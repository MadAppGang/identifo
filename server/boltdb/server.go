package boltdb

import (
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/plugin/shared"
	"github.com/madappgang/identifo/server"
)

// NewServer creates new backend service with BoltDB support.
func NewServer(settings model.ServerSettings, cors *model.CorsOptions, plugins shared.Plugins, serverOptions ...func(*server.Server) error) (model.Server, error) {
	dbComposer, err := NewComposer(settings, plugins)
	if err != nil {
		return nil, err
	}
	return server.NewServer(settings, dbComposer, nil, cors, serverOptions...)
}
