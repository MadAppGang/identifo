package dynamodb

import (
	"fmt"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/server"
)

// NewServer creates new backend service with DynamoDB support.
func NewServer(settings model.ServerSettings, options ...func(*DatabaseComposer) error) (model.Server, error) {
	if settings.Database.DBType != "dynamodb" {
		return nil, fmt.Errorf("Incorrect database type %s for DynamoDB-backed server", settings.Database.DBType)
	}

	dbComposer, err := NewComposer(settings, options...)
	if err != nil {
		return nil, err
	}
	return server.NewServer(settings, dbComposer)
}
