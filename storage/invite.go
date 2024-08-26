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

// NewInviteStorage creates new invite storage from settings
func NewInviteStorage(
	logger *slog.Logger,
	settings model.DatabaseSettings) (model.InviteStorage, error) {
	switch settings.Type {
	case model.DBTypeBoltDB:
		return boltdb.NewInviteStorage(logger, settings.BoltDB)
	case model.DBTypeMongoDB:
		return mongo.NewInviteStorage(logger, settings.Mongo)
	case model.DBTypeDynamoDB:
		return dynamodb.NewInviteStorage(logger, settings.Dynamo)
	case model.DBTypeFake:
		fallthrough
	case model.DBTypeMem:
		return mem.NewInviteStorage()
	default:
		return nil, fmt.Errorf("invite storage type is not supported %s ", settings.Type)
	}
}
