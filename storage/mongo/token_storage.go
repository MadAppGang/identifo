package mongo

import (
	"github.com/madappgang/identifo/model"
	"gopkg.in/mgo.v2/bson"
)

const (
	// TokensCollection is a collection to store refresh tokens.
	TokensCollection = "RefreshTokens"
)

// NewTokenStorage creates a MongoDB token storage.
func NewTokenStorage(db *DB) (model.TokenStorage, error) {
	return &TokenStorage{db: db}, nil
}

// TokenStorage is a MongoDB token storage.
type TokenStorage struct {
	db *DB
}

// SaveToken saves token in the database.
func (ts *TokenStorage) SaveToken(token string) error {
	if len(token) == 0 {
		return model.ErrorWrongDataFormat
	}
	s := ts.db.Session(TokensCollection)
	defer s.Close()

	var t = Token{Token: token, ID: bson.NewObjectId()}
	err := s.C.Insert(t)
	return err
}

// HasToken returns true if the token is present in the storage.
func (ts *TokenStorage) HasToken(token string) bool {
	s := ts.db.Session(TokensCollection)
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

// DeleteToken removes token from the storage.
func (ts *TokenStorage) DeleteToken(token string) error {
	s := ts.db.Session(TokensCollection)
	defer s.Close()

	if _, err := s.C.RemoveAll(bson.M{"token": token}); err != nil {
		return err
	}
	return nil
}

// Close closes database connection.
func (ts *TokenStorage) Close() {
	ts.db.Close()
}

// Token is struct to store tokens in database.
type Token struct {
	ID    bson.ObjectId `bson:"_id,omitempty"`
	Token string        `bson:"token,omitempty"`
}
