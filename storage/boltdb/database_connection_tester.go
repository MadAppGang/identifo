package boltdb

import (
	"fmt"

	"github.com/madappgang/identifo/v2/model"
	bolt "go.etcd.io/bbolt"
)

type ConnectionTester struct {
	settings model.BoltDBDatabaseSettings
}

// NewConnectionTester creates a BoltDB connection tester

func NewConnectionTester(settings model.BoltDBDatabaseSettings) model.ConnectionTester {
	return &ConnectionTester{settings: settings}
}

func (ct *ConnectionTester) Connect() error {
	if len(ct.settings.Path) == 0 {
		return ErrorEmptyDatabasePath
	}

	db, err := InitDB(ct.settings.Path)
	if err != nil {
		return err
	}

	// trying to create test bucket
	if err := db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte("TestConnection")); err != nil {
			return fmt.Errorf("error creating bucket: %s", err)
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}
