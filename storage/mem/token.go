package mem

import (
	"github.com/madappgang/identifo/model"
)

// NewTokenStorage creates an in-memory token storage.
func NewTokenStorage() (model.TokenStorage, error) {
	return &TokenStorage{storage: make(map[string]bool)}, nil
}

// TokenStorage is an in-memory token storage.
// Please do not use it in production, it has no disk swap or persistent cache support.
type TokenStorage struct {
	storage map[string]bool
}

// SaveToken saves token in memory.
func (ts *TokenStorage) SaveToken(token string) error {
	ts.storage[token] = true
	return nil
}

// HasToken returns true if the token is present in the storage.
func (ts *TokenStorage) HasToken(token string) bool {
	has := ts.storage[token]
	return has
}

// DeleteToken removes token from memory storage.
// Actually, just marks it as deleted.
func (ts *TokenStorage) DeleteToken(token string) error {
	ts.storage[token] = false
	return nil
}

// Close clears storage.
func (ts *TokenStorage) Close() {
	for k := range ts.storage {
		delete(ts.storage, k)
	}
}
