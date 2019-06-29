package boltdb

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
	"github.com/madappgang/identifo/model"
)

const (
	// TokenBucket is a name for bucket with tokens.
	TokenBucket = "Tokens"
)

// NewTokenStorage creates a BoltDB token storage.
func NewTokenStorage(db *bolt.DB) (model.TokenStorage, error) {
	ts := &TokenStorage{db: db}
	// ensure that we have needed bucket in the database
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
	db *bolt.DB
}

// SaveToken saves token in the storage.
func (ts *TokenStorage) SaveToken(token string) error {
	return ts.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TokenBucket))
		//we use token as key and value
		return b.Put([]byte(token), []byte(token))
	})
}

// HasToken returns true if the token in present in the storage.
func (ts *TokenStorage) HasToken(token string) bool {
	var res bool
	if err := ts.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TokenBucket))
		//we use token as key and value
		res = b.Get([]byte(token)) != nil
		return nil
	}); err != nil {
		return false
	}

	return res
}

// RevokeToken removes token from the storage.
func (ts *TokenStorage) RevokeToken(token string) error {
	return ts.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TokenBucket))
		//we use token as key and value
		return b.Delete([]byte(token))
	})
}

// Close closes underlying database.
func (ts *TokenStorage) Close() {
	if err := ts.db.Close(); err != nil {
		log.Printf("Error closing token storage: %s\n", err)
	}
}
