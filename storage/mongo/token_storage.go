package mongo

import (
	"context"
	"time"

	"github.com/madappgang/identifo/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const tokensCollectionName = "RefreshTokens"

// NewTokenStorage creates a MongoDB token storage.
func NewTokenStorage(settings model.MongodDatabaseSettings) (model.TokenStorage, error) {
	if len(settings.ConnectionString) == 0 || len(settings.DatabaseName) == 0 {
		return nil, ErrorEmptyConnectionStringDatabase
	}

	// create database
	db, err := NewDB(settings.ConnectionString, settings.DatabaseName)
	if err != nil {
		return nil, err
	}

	coll := db.Database.Collection(tokensCollectionName)
	return &TokenStorage{coll: coll, timeout: 30 * time.Second}, nil
}

// TokenStorage is a MongoDB token storage.
type TokenStorage struct {
	coll    *mongo.Collection
	timeout time.Duration
}

// SaveToken saves token in the database.
func (ts *TokenStorage) SaveToken(token string) error {
	if len(token) == 0 {
		return model.ErrorWrongDataFormat
	}

	ctx, cancel := context.WithTimeout(context.Background(), ts.timeout)
	defer cancel()

	t := Token{Token: token, ID: primitive.NewObjectID().Hex()}
	_, err := ts.coll.InsertOne(ctx, t)
	return err
}

// HasToken returns true if the token is present in the storage.
func (ts *TokenStorage) HasToken(token string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), ts.timeout)
	defer cancel()

	var t Token
	if err := ts.coll.FindOne(ctx, bson.M{"token": token}).Decode(&t); err != nil {
		return false
	}
	if t.Token == token {
		return true
	}
	return false
}

// DeleteToken removes token from the storage.
func (ts *TokenStorage) DeleteToken(token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), ts.timeout)
	defer cancel()

	_, err := ts.coll.DeleteMany(ctx, bson.M{"token": token})
	return err
}

// Close is a no-op.
func (ts *TokenStorage) Close() {}

// Token is struct to store tokens in database.
type Token struct {
	ID    string `bson:"_id,omitempty"` // TODO: Make use of jti claim.
	Token string `bson:"token,omitempty"`
}
