package storage

import (
	"fmt"
	"log/slog"

	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage/boltdb"
	"github.com/madappgang/identifo/v2/storage/dynamodb"
	"github.com/madappgang/identifo/v2/storage/mem"
	"github.com/madappgang/identifo/v2/storage/mongo"
)

// NewAppStorage creates new app storage from settings
func NewAppStorage(
	logger *slog.Logger,
	settings model.DatabaseSettings) (model.AppStorage, error) {
	switch settings.Type {
	case model.DBTypeBoltDB:
		return boltdb.NewAppStorage(logger, settings.BoltDB)
	case model.DBTypeMongoDB:
		return mongo.NewAppStorage(logger, settings.Mongo)
	case model.DBTypeDynamoDB:
		return dynamodb.NewAppStorage(logger, settings.Dynamo)
	case model.DBTypeFake:
		fallthrough
	case model.DBTypeMem:
		return mem.NewAppStorage(logger)
	default:
		return nil, fmt.Errorf("app storage type is not supported")
	}
}
