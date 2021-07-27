package storage

import (
	"fmt"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/storage/boltdb"
	"github.com/madappgang/identifo/storage/dynamodb"
	"github.com/madappgang/identifo/storage/mem"
	"github.com/madappgang/identifo/storage/mongo"
)

// NewUserStorage creates new users storage from settings
func NewUserStorage(settings model.DatabaseSettings) (model.UserStorage, error) {
	switch settings.Type {
	case model.DBTypeBoltDB:
		return boltdb.NewUserStorage(settings.BoltDB)
	case model.DBTypeMongoDB:
		return mongo.NewUserStorage(settings.Mongo)
	case model.DBTypeDynamoDB:
		return dynamodb.NewUserStorage(settings.Dynamo)
	case model.DBTypeFake:
		return mem.NewUserStorage()
	default:
		return nil, fmt.Errorf("user storage type is not supported %s ", settings.Type)
	}
}
