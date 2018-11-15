package mem

import (
	"github.com/madappgang/identifo/model"
)

//NewTokenStorage created in memory token sotrage
func NewTokenStorage() model.TokenStorage {
	ts := TokenStorage{}
	ts.storage = make(map[string]bool)
	return &ts
}

//TokenStorage im memory token storage
//please don't use it in production
//no disk swap support
//no persisten cache support
type TokenStorage struct {
	storage map[string]bool
}

//SaveToken save token in memory
func (ts *TokenStorage) SaveToken(token string) error {
	ts.storage[token] = true
	return nil
}

//HasToken returns if the token in the storage
func (ts *TokenStorage) HasToken(token string) bool {
	has := ts.storage[token]
	return has
}

//RevokeToken removes token from memory storage
//actually just mark the token as deleted
func (ts *TokenStorage) RevokeToken(token string) error {
	ts.storage[token] = false
	return nil
}
