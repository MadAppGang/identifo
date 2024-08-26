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

// NewTokenBlacklistStorage creates new tokens blacklist storage from settings
func NewTokenBlacklistStorage(
	logger *slog.Logger,
	settings model.DatabaseSettings,
) (model.TokenBlacklist, error) {
	switch settings.Type {
	case model.DBTypeBoltDB:
		return boltdb.NewTokenBlacklist(logger, settings.BoltDB)
	case model.DBTypeMongoDB:
		return mongo.NewTokenBlacklist(logger, settings.Mongo)
	case model.DBTypeDynamoDB:
		return dynamodb.NewTokenBlacklist(logger, settings.Dynamo)
	case model.DBTypeFake:
		fallthrough
	case model.DBTypeMem:
		return mem.NewTokenBlacklist()
	default:
		return nil, fmt.Errorf("token blacklist storage type is not supported %s ", settings.Type)
	}
}
