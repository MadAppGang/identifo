package model

import (
	"context"
	"time"
)

// TokenStorage is a storage for tokens.
type TokenStorage interface {
	SaveToken(ctx context.Context, token TokenStorageEntity) error
	TokenByID(ctx context.Context, id string) (TokenStorageEntity, error)
	TokenByRaw(ctx context.Context, raw string) (TokenStorageEntity, error)
	DeleteToken(ctx context.Context, id string) error
	Close()
}

// TokenStorageEntity keep information about tokens in the system
type TokenStorageEntity struct {
	ID        string              `json:"id" bson:"_id"`
	RawToken  string              `json:"raw_token" bson:"raw_token"`
	TokenType TokenType           `json:"token_type" bson:"token_type"`
	AddedAt   time.Time           `json:"added_at" bson:"added_at"`
	AddedBy   TokenStorageAddedBy `json:"added_by" bson:"added_by"`
	Comments  string              `json:"comments" bson:"comments"`
}

type TokenStorageAddedBy string

var (
	TokenStorageAddedByUser       TokenStorageAddedBy = "user"       // user added it by refreshing the token or so
	TokenStorageAddedByAdmin      TokenStorageAddedBy = "admin"      // admin panel
	TokenStorageAddedByManagement TokenStorageAddedBy = "management" // management api
)

// JWTKeys are keys used for signing and verifying JSON web tokens.
type JWTKeys struct {
	Public  interface{}
	Private interface{}
}
