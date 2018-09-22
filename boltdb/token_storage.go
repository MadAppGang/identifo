package boltdb

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/madappgang/identifo/model"
)

const (
	//TokenBucket bucket name with refresh tokens
	TokenBucket = "Tokens"
)

//NewTokenStorage created in embedded token sotrage
func NewTokenStorage(db *bolt.DB) (model.TokenStorage, error) {
	ts := TokenStorage{}
	ts.db = db
	//ensure we have app's bucket in the database
	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(TokenBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &ts, nil

}

//TokenStorage im embedded token storage
type TokenStorage struct {
	db *bolt.DB
}

//SaveToken save token in database
func (ts *TokenStorage) SaveToken(token string) error {
	return ts.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TokenBucket))
		//we use token as key and value
		return b.Put([]byte(token), []byte(token))
	})
}

//HasToken returns true if the token in the storage
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

//RevokeToken removes token from the storage
func (ts *TokenStorage) RevokeToken(token string) error {
	return ts.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TokenBucket))
		//we use token as key and value
		return b.Delete([]byte(token))
	})
}
