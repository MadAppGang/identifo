package storage

import (
	"fmt"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/storage/boltdb"
	"github.com/madappgang/identifo/storage/dynamodb"
	"github.com/madappgang/identifo/storage/mem"
	"github.com/madappgang/identifo/storage/mongo"
)

// NewAppStorage creates new app storage from settings
func NewAppStorage(settings model.DatabaseSettings) (model.AppStorage, error) {
	switch settings.Type {
	case model.DBTypeBoltDB:
		return boltdb.NewAppStorage(settings.BoltDB)
	case model.DBTypeMongoDB:
		return mongo.NewAppStorage(settings.Mongo)
	case model.DBTypeDynamoDB:
		return dynamodb.NewAppStorage(settings.Dynamo)
	case model.DBTypeFake:
		return mem.NewAppStorage()
	default:
		return nil, fmt.Errorf("App storage type is not supported")
	}
}
