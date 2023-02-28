package storage

import (
	"fmt"

	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage/boltdb"
	"github.com/madappgang/identifo/v2/storage/dynamodb"
	"github.com/madappgang/identifo/v2/storage/mem"
	"github.com/madappgang/identifo/v2/storage/mongo"
)

// NewManagementKeys creates new management keys storage from settings.
func NewManagementKeys(settings model.DatabaseSettings) (model.ManagementKeysStorage, error) {
	switch settings.Type {
	case model.DBTypeBoltDB:
		return boltdb.NewManagementKeysStorage(settings.BoltDB)
	case model.DBTypeMongoDB:
		return mongo.NewManagementKeysStorage(settings.Mongo)
	case model.DBTypeDynamoDB:
		return dynamodb.NewManagementKeysStorage(settings.Dynamo)
	case model.DBTypeFake:
		fallthrough
	case model.DBTypeMem:
		return mem.NewManagementKeysStorage()
	default:
		return nil, fmt.Errorf("token storage type is not supported %s ", settings.Type)
	}
}
