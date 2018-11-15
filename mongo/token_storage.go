package mongo

import (
	"github.com/madappgang/identifo/model"
	"gopkg.in/mgo.v2/bson"
)

const (
	//TokensCollection is collection to store refresh tokens
	TokensCollection = "RefreshTokens"
)

//NewTokenStorage created in embedded token sotrage
func NewTokenStorage(db *DB) (model.TokenStorage, error) {
	ts := TokenStorage{db: db}
	return &ts, nil
}

//TokenStorage omnogbased token storage
type TokenStorage struct {
	db *DB
}

//SaveToken save token in database
func (ts *TokenStorage) SaveToken(token string) error {
	if len(token) == 0 {
		return model.ErrorWrongDataFormat
	}
	s := ts.db.Session(TokensCollection)
	defer s.Close()

	var t Token
	t.Token = token
	t.ID = bson.NewObjectId()
	if err := s.C.Insert(t); err != nil {
		return err
	}
	return nil
}

//HasToken returns true if the token in the storage
func (ts *TokenStorage) HasToken(token string) bool {
	s := ts.db.Session(TokensCollection)
	defer s.Close()

	var t Token
	q := bson.M{"token": token}
	if err := s.C.Find(q).One(&t); err != nil {
		return false
	}
	if t.Token == token {
		return true
	}

	return false
}

//RevokeToken removes token from the storage
func (ts *TokenStorage) RevokeToken(token string) error {
	s := ts.db.Session(TokensCollection)
	defer s.Close()

	q := bson.M{"token": token}
	if _, err := s.C.RemoveAll(q); err != nil {
		return err
	}
	return nil

}

//Token is struct to store tokens in database
type Token struct {
	ID    bson.ObjectId `bson:"_id,omitempty"`
	Token string        `bson:"token,omitempty"`
}
