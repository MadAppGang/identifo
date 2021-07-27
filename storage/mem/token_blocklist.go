package mem

import (
	"github.com/madappgang/identifo/model"
)

// NewTokenBlacklist creates an in-memory token storage.
func NewTokenBlacklist() (model.TokenBlacklist, error) {
	return &TokenBlacklist{storage: make(map[string]bool)}, nil
}

// TokenBlacklist is an in-memory token storage.
// Please do not use it in production, it has no disk swap or persistent cache support.
type TokenBlacklist struct {
	storage map[string]bool
}

// Add blacklists token.
func (tb *TokenBlacklist) Add(token string) error {
	tb.storage[token] = true
	return nil
}

// IsBlacklisted returns true if the token is blacklisted.
func (tb *TokenBlacklist) IsBlacklisted(token string) bool {
	has := tb.storage[token]
	return has
}

// Close clears storage.
func (tb *TokenBlacklist) Close() {
	for k := range tb.storage {
		delete(tb.storage, k)
	}
}
