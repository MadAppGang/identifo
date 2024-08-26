package boltdb

import (
	"fmt"
	"log/slog"

	"github.com/madappgang/identifo/v2/logging"
	"github.com/madappgang/identifo/v2/model"
	bolt "go.etcd.io/bbolt"
)

const (
	// TokenBucket is a name for bucket with tokens.
	TokenBucket = "Tokens"
)

// NewTokenStorage creates a BoltDB token storage.
func NewTokenStorage(
	logger *slog.Logger,
	settings model.BoltDBDatabaseSettings,
) (model.TokenStorage, error) {
	if len(settings.Path) == 0 {
		return nil, ErrorEmptyDatabasePath
	}

	// init database
	db, err := InitDB(settings.Path)
	if err != nil {
		return nil, err
	}

	ts := &TokenStorage{
		logger: logger,
		db:     db,
	}
	// Ensure that we have needed bucket in the database.
	if err := db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(TokenBucket)); err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return ts, nil
}

// TokenStorage is a BoltDB token storage.
type TokenStorage struct {
	logger *slog.Logger
	db     *bolt.DB
}

// SaveToken saves token in the storage.
func (ts *TokenStorage) SaveToken(token string) error {
	return ts.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TokenBucket))
		// We use token as key and value.
		return b.Put([]byte(token), []byte(token))
	})
}

// HasToken returns true if the token is present in the storage.
func (ts *TokenStorage) HasToken(token string) bool {
	var res bool
	if err := ts.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TokenBucket))
		// We use token as key and value.
		res = b.Get([]byte(token)) != nil
		return nil
	}); err != nil {
		return false
	}

	return res
}

// DeleteToken removes token from the storage.
func (ts *TokenStorage) DeleteToken(token string) error {
	return ts.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TokenBucket))
		// We use token as key and value.
		return b.Delete([]byte(token))
	})
}

// Close closes underlying database.
func (ts *TokenStorage) Close() {
	if err := CloseDB(ts.db); err != nil {
		ts.logger.Error("Error closing token storage", logging.FieldError, err)
	}
}
