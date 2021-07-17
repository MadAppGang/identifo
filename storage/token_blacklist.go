package storage

import (
	"fmt"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/storage/boltdb"
	"github.com/madappgang/identifo/storage/dynamodb"
	"github.com/madappgang/identifo/storage/mem"
	"github.com/madappgang/identifo/storage/mongo"
)

// NewTokenBlacklistStorage creates new tokens blacklist storage from settings
func NewTokenBlacklistStorage(settings model.DatabaseSettings) (model.TokenBlacklist, error) {
	switch settings.Type {
	case model.DBTypeBoltDB:
		return boltdb.NewTokenBlacklist(settings.BoltDB)
	case model.DBTypeMongoDB:
		return mongo.NewTokenBlacklist(settings.Mongo)
	case model.DBTypeDynamoDB:
		return dynamodb.NewTokenBlacklist(settings.Dynamo)
	case model.DBTypeFake:
		return mem.NewTokenBlacklist()
	default:
		return nil, fmt.Errorf("token blacklist storage type is not supported %s ", settings.Type)
	}
}
