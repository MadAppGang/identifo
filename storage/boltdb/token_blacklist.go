package boltdb

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
	"github.com/madappgang/identifo/model"
)

const (
	// BlacklistedTokenBucket is a name for bucket with tokens blacklist.
	BlacklistedTokenBucket = "BlacklistedTokens"
)

// NewTokenBlacklist creates a token blacklist in BoltDB.
func NewTokenBlacklist(db *bolt.DB) (model.TokenBlacklist, error) {
	tb := &TokenBlacklist{db: db}
	if err := tb.db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(BlacklistedTokenBucket)); err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return tb, nil
}

// TokenBlacklist is a BoltDB token blacklist.
type TokenBlacklist struct {
	db *bolt.DB
}

// Add adds token in the blacklist.
func (tb *TokenBlacklist) Add(token string) error {
	return tb.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlacklistedTokenBucket))
		// We use token as key and value.
		return b.Put([]byte(token), []byte(token))
	})
}

// IsBlacklisted returns true if the token is blacklisted.
func (tb *TokenBlacklist) IsBlacklisted(token string) bool {
	var res bool
	if err := tb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlacklistedTokenBucket))
		// We use token as key and value.
		res = b.Get([]byte(token)) != nil
		return nil
	}); err != nil {
		return false
	}
	return res
}

// Close closes underlying database.
func (tb *TokenBlacklist) Close() {
	if err := tb.db.Close(); err != nil {
		log.Printf("Error closing token blacklist storage: %s\n", err)
	}
}
