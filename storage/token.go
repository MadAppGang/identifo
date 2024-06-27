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

// NewTokenStorage creates new tokens storage from settings
func NewTokenStorage(
	logger *slog.Logger,
	settings model.DatabaseSettings) (model.TokenStorage, error) {
	switch settings.Type {
	case model.DBTypeBoltDB:
		return boltdb.NewTokenStorage(logger, settings.BoltDB)
	case model.DBTypeMongoDB:
		return mongo.NewTokenStorage(logger, settings.Mongo)
	case model.DBTypeDynamoDB:
		return dynamodb.NewTokenStorage(logger, settings.Dynamo)
	case model.DBTypeFake:
		fallthrough
	case model.DBTypeMem:
		return mem.NewTokenStorage()
	default:
		return nil, fmt.Errorf("token storage type is not supported %s ", settings.Type)
	}
}
