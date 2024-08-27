package storage

import (
	"fmt"
	"log/slog"

	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage/boltdb"
	"github.com/madappgang/identifo/v2/storage/dynamodb"
	"github.com/madappgang/identifo/v2/storage/grpc"
	"github.com/madappgang/identifo/v2/storage/mem"
	"github.com/madappgang/identifo/v2/storage/mongo"
	"github.com/madappgang/identifo/v2/storage/plugin"
)

// NewUserStorage creates new users storage from settings
func NewUserStorage(
	logger *slog.Logger,
	settings model.DatabaseSettings) (model.UserStorage, error) {
	switch settings.Type {
	case model.DBTypeBoltDB:
		return boltdb.NewUserStorage(logger, settings.BoltDB)
	case model.DBTypeMongoDB:
		return mongo.NewUserStorage(logger, settings.Mongo)
	case model.DBTypeDynamoDB:
		return dynamodb.NewUserStorage(logger, settings.Dynamo)
	case model.DBTypeFake:
		fallthrough
	case model.DBTypeMem:
		return mem.NewUserStorage()
	case model.DBTypePlugin:
		return plugin.NewUserStorage(logger, settings.Plugin)
	case model.DBTypeGRPC:
		return grpc.NewUserStorage(settings.GRPC)
	default:
		return nil, fmt.Errorf("user storage type is not supported %s ", settings.Type)
	}
}
