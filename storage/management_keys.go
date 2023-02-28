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
		return boltdb.NewManagementKeys(settings.BoltDB)
	case model.DBTypeMongoDB:
		return mongo.NewManagementKeys(settings.Mongo)
	case model.DBTypeDynamoDB:
		return dynamodb.NewManagementKeys(settings.Dynamo)
	case model.DBTypeFake:
		fallthrough
	case model.DBTypeMem:
		return mem.NewManagementKeys()
	default:
		return nil, fmt.Errorf("token storage type is not supported %s ", settings.Type)
	}
}
