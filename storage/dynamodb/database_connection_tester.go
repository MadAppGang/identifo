package dynamodb

import (
	"github.com/madappgang/identifo/model"
)

const testTableName = "test_connection_table"

type ConnectionTester struct {
	settings model.DynamoDatabaseSettings
}

// NewConnectionTester creates a BoltDB connection tester

func NewConnectionTester(settings model.DynamoDatabaseSettings) model.ConnectionTester {
	return &ConnectionTester{settings: settings}
}

func (ct *ConnectionTester) Connect() error {
	if len(ct.settings.Endpoint) == 0 || len(ct.settings.Region) == 0 {
		return ErrorEmptyEndpointRegion
	}

	// create database
	_, err := NewDB(ct.settings.Endpoint, ct.settings.Region)
	if err != nil {
		return err
	}

	return nil
}
