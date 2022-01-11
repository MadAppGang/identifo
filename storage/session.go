package storage

import (
	"fmt"

	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage/dynamodb"
	"github.com/madappgang/identifo/v2/storage/mem"
	"github.com/madappgang/identifo/v2/storage/redis"
)

// NewSessionStorage creates new sessions storage from settings
func NewSessionStorage(settings model.SessionStorageSettings) (model.SessionStorage, error) {
	switch settings.Type {
	case model.SessionStorageRedis:
		return redis.NewSessionStorage(settings.Redis)
	case model.SessionStorageMem:
		return mem.NewSessionStorage(), nil
	case model.SessionStorageDynamoDB:
		return dynamodb.NewSessionStorage(settings.Dynamo)
	default:
		return nil, fmt.Errorf("session storage type is not supported")
	}
}
