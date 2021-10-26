package storage

import (
	"fmt"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/storage/boltdb"
	"github.com/madappgang/identifo/storage/dynamodb"
	"github.com/madappgang/identifo/storage/fs"
	"github.com/madappgang/identifo/storage/mem"
	"github.com/madappgang/identifo/storage/mongo"
	"github.com/madappgang/identifo/storage/s3"
)

func NewConnectionTester(settings model.TestConnection) (model.ConnectionTester, error) {
	switch settings.Type {
	case model.TTDatabase:
		if settings.Database == nil {
			return nil, fmt.Errorf("empty connection settings for testing type %s", settings.Type)
		}
		return NewDatabaseConnectionTester(*settings.Database), nil
	case model.TTKeyStorage:
		if settings.KeyStorage == nil {
			return nil, fmt.Errorf("empty connection settings for testing type %s", settings.Type)
		}
		return NewKeyStorageConnectionTester(*settings.KeyStorage), nil
	case model.TTFileStorage:
		if settings.FileStorage == nil {
			return nil, fmt.Errorf("empty file storage settings for testing type %s", settings.Type)
		}
		return NewFileStorageConnectionTester(*settings.FileStorage), nil
	}

	return nil, fmt.Errorf("unknown settings type fro testing: %v", settings.Type)
}

func NewDatabaseConnectionTester(settings model.DatabaseSettings) model.ConnectionTester {
	switch settings.Type {
	case model.DBTypeBoltDB:
		return boltdb.NewConnectionTester(settings.BoltDB)
	case model.DBTypeMongoDB:
		return mongo.NewConnectionTester(settings.Mongo)
	case model.DBTypeDynamoDB:
		return dynamodb.NewConnectionTester(settings.Dynamo)
	case model.DBTypeFake:
		return mem.NewConnectionTester()
	}
	return nil
}

func NewKeyStorageConnectionTester(settings model.KeyStorageSettings) model.ConnectionTester {
	switch settings.Type {
	case model.KeyStorageTypeLocal:
		return fs.NewKeyStorageConnectionTester(settings.File)
	case model.KeyStorageTypeS3:
		return s3.NewKeyStorageConnectionTester(settings.S3)
	}
	return nil
}

func NewFileStorageConnectionTester(settings model.FileStorageTestSettings) model.ConnectionTester {
	switch settings.FileStorage.Type {
	case model.FileStorageTypeNone, model.FileStorageTypeDefault:
		return AlwaysHappyConnectionTester{}
	case model.FileStorageTypeLocal:
		return fs.NewFSConnectionTester(settings.FileStorage.Local, settings.ExpectedFiles)
	}
	return nil
}

type AlwaysFailedConnectionTester struct{}

func (ct AlwaysFailedConnectionTester) Connect() error {
	return fmt.Errorf("unsupported connection type")
}

type AlwaysHappyConnectionTester struct{}

func (ct AlwaysHappyConnectionTester) Connect() error {
	return nil
}
