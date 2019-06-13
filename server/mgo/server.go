package mgo

import (
	"fmt"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/server"
)

func init() {

}

// NewServer creates new backend service with MongoDB support.
func NewServer(settings model.ServerSettings, options ...func(*DatabaseComposer) error) (model.Server, error) {
	if settings.DBType != "mongodb" {
		return nil, fmt.Errorf("Incorrect database type %s for MongoDB-backed server", settings.DBType)
	}

	dbComposer, err := NewComposer(settings, options...)
	if err != nil {
		return nil, err
	}
	return server.NewServer(settings, dbComposer)
}
