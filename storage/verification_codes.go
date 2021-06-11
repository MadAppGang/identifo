package storage

import (
	"fmt"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/storage/boltdb"
	"github.com/madappgang/identifo/storage/dynamodb"
	"github.com/madappgang/identifo/storage/mem"
	"github.com/madappgang/identifo/storage/mongo"
)

// NewVerificationCodesStorage creates new verification codes storage from settings
func NewVerificationCodesStorage(settings model.DatabaseSettings) (model.VerificationCodeStorage, error) {
	switch settings.Type {
	case model.DBTypeBoltDB:
		return boltdb.NewVerificationCodeStorage(settings.BoltDB)
	case model.DBTypeMongoDB:
		return mongo.NewVerificationCodeStorage(settings.Mongo)
	case model.DBTypeDynamoDB:
		return dynamodb.NewVerificationCodeStorage(settings.Dynamo)
	case model.DBTypeFake:
		return mem.NewVerificationCodeStorage()
	default:
		return nil, fmt.Errorf("verification code storage type is not supported %s ", settings.Type)
	}
}
