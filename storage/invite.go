package storage

import (
	"fmt"

	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage/boltdb"
	"github.com/madappgang/identifo/v2/storage/dynamodb"
	"github.com/madappgang/identifo/v2/storage/mongo"
)

// NewInviteStorage creates new invite storage from settings
func NewInviteStorage(settings model.DatabaseSettings) (model.InviteStorage, error) {
	switch settings.Type {
	case model.DBTypeBoltDB:
		return boltdb.NewInviteStorage(settings.BoltDB)
	case model.DBTypeMongoDB:
		return mongo.NewInviteStorage(settings.Mongo)
	case model.DBTypeDynamoDB:
		return dynamodb.NewInviteStorage(settings.Dynamo)
	case model.DBTypeFake:
		fallthrough
	default:
		return nil, fmt.Errorf("invite storage type is not supported %s ", settings.Type)
	}
}
