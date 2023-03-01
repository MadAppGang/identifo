package model

import (
	"context"
	"time"
)

// ManagementKeysStorage storage to persist management keys
type ManagementKeysStorage interface {
	GetKey(ctx context.Context, id string) (ManagementKey, error)
	CreateKey(ctx context.Context, name string, scopes []string) (ManagementKey, error)
	DisableKey(ctx context.Context, id string) (ManagementKey, error)
	RenameKey(ctx context.Context, id, name string) (ManagementKey, error)
	ChangeScopesForKey(ctx context.Context, id string, scopes []string) (ManagementKey, error)
	UseKey(ctx context.Context, id string) (ManagementKey, error)
	ImportJSON(data []byte, clearOldData bool) error

	GeyAllKeys(ctx context.Context) ([]ManagementKey, error)
}

// ManagementKey secret management key to communicate with management api
type ManagementKey struct {
	ID        string     `json:"id" bson:"_id"`
	Secret    string     `json:"secret" bson:"secret"`
	Name      string     `json:"name" bson:"name"`
	Active    bool       `json:"active" bson:"active"`
	Scopes    []string   `json:"scopes" bson:"scopes"`
	CreatedAt time.Time  `json:"created_at" bson:"createdAt"`
	LastUsed  time.Time  `json:"last_used" bson:"lastUsed"`
	ValidTill *time.Time `json:"valid_till" bson:"validTill"`
}
