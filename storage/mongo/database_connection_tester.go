package mongo

import (
	"github.com/madappgang/identifo/v2/model"
)

type ConnectionTester struct {
	settings model.MongoDatabaseSettings
}

// NewConnectionTester creates a BoltDB connection tester

func NewConnectionTester(settings model.MongoDatabaseSettings) model.ConnectionTester {
	return &ConnectionTester{settings: settings}
}

func (ct *ConnectionTester) Connect() error {
	if len(ct.settings.ConnectionString) == 0 || len(ct.settings.DatabaseName) == 0 {
		return ErrorEmptyConnectionStringDatabase
	}

	// create or connect to database
	_, err := NewDB(ct.settings.ConnectionString, ct.settings.DatabaseName)
	if err != nil {
		return err
	}

	return nil
}
