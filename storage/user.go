package storage

import (
	"fmt"

	"github.com/madappgang/identifo/v2/model"
)

// NewUserStorage creates new users storage from settings
func NewUserStorage(settings model.DatabaseSettings) (model.UserStorage, error) {
	// TODO: Implement storages and active here
	// ! impossible to go next without it
	// ! let's implement boltdb first
	switch settings.Type {
	// case model.DBTypeBoltDB:
	// 	return boltdb.NewUserStorage(settings.BoltDB)
	// case model.DBTypeMongoDB:
	// 	return mongo.NewUserStorage(settings.Mongo)
	// case model.DBTypeDynamoDB:
	// 	return dynamodb.NewUserStorage(settings.Dynamo)
	// case model.DBTypePlugin:
	// 	return plugin.NewUserStorage(settings.Plugin)
	// case model.DBTypeGRPC:
	// 	return grpc.NewUserStorage(settings.GRPC)
	default:
		return nil, fmt.Errorf("user storage type is not supported %s ", settings.Type)
	}
}
