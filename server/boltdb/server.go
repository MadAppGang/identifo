package boltdb

import (
	"fmt"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/server"
)

// NewServer creates new backend service with BoltDB support.
func NewServer(settings model.ServerSettings, options ...func(*DatabaseComposer) error) (model.Server, error) {
	if settings.DBType != "boltdb" {
		return nil, fmt.Errorf("Incorrect database type %s for BoltDB-backed server", settings.DBType)
	}

	dbComposer, err := NewComposer(settings, options...)
	if err != nil {
		return nil, err
	}
	return server.NewServer(settings, dbComposer)
}
