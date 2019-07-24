package mongo

import (
	"github.com/madappgang/identifo/model"
	"gopkg.in/mgo.v2/bson"
)

const (
	// BlacklistedTokensCollection is a collection where blacklisted tokens are stored.
	BlacklistedTokensCollection = "BlacklistedTokens"
)

// NewTokenBlacklist creates a MongoDB token storage.
func NewTokenBlacklist(db *DB) (model.TokenBlacklist, error) {
	return &TokenBlacklist{db: db}, nil
}

// TokenBlacklist is a MongoDB token blacklist.
type TokenBlacklist struct {
	db *DB
}

// Add adds token to the blacklist.
func (tb *TokenBlacklist) Add(token string) error {
	if len(token) == 0 {
		return model.ErrorWrongDataFormat
	}
	s := tb.db.Session(BlacklistedTokensCollection)
	defer s.Close()

	var t = Token{Token: token, ID: bson.NewObjectId()}
	err := s.C.Insert(t)
	return err
}

// IsBlacklisted returns true if the token is present in the blacklist.
func (tb *TokenBlacklist) IsBlacklisted(token string) bool {
	s := tb.db.Session(BlacklistedTokensCollection)
	defer s.Close()

	var t Token
	if err := s.C.Find(bson.M{"token": token}).One(&t); err != nil {
		return false
	}
	if t.Token == token {
		return true
	}
	return false
}

// Close closes database connection.
func (tb *TokenBlacklist) Close() {
	tb.db.Close()
}
