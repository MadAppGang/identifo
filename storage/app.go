package storage

import (
	"fmt"

	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage/boltdb"
	"github.com/madappgang/identifo/v2/storage/dynamodb"
	"github.com/madappgang/identifo/v2/storage/mem"
	"github.com/madappgang/identifo/v2/storage/mongo"
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
		fallthrough
	case model.DBTypeMem:
		return mem.NewAppStorage()
	default:
		return nil, fmt.Errorf("App storage type is not supported")
	}
}
