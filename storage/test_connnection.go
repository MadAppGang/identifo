package storage

import (
	"fmt"

	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage/boltdb"
	"github.com/madappgang/identifo/v2/storage/dynamodb"
	"github.com/madappgang/identifo/v2/storage/fs"
	"github.com/madappgang/identifo/v2/storage/mongo"
	"github.com/madappgang/identifo/v2/storage/s3"
)

var spaFileStorageExpectedFiles = []string{"index.html"}

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
	case model.TTSPAFileStorage:
		if settings.FileStorage == nil {
			return nil, fmt.Errorf("empty file storage settings for testing type %s", settings.Type)
		}
		return NewFileStorageConnectionTester(*settings.FileStorage, spaFileStorageExpectedFiles), nil
	case model.TTEmailTemplateStorage:
		if settings.FileStorage == nil {
			return nil, fmt.Errorf("empty file storage settings for testing type %s", settings.Type)
		}
		return NewFileStorageConnectionTester(*settings.FileStorage, model.AllEmailTemplatesFileNames()), nil
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
	}
	return nil
}

func NewKeyStorageConnectionTester(settings model.FileStorageSettings) model.ConnectionTester {
	switch settings.Type {
	case model.FileStorageTypeLocal:
		return fs.NewKeyStorageConnectionTester(settings.Local)
	case model.FileStorageTypeS3:
		return s3.NewKeyStorageConnectionTester(settings.S3)
	}
	return nil
}

func NewFileStorageConnectionTester(settings model.FileStorageSettings, expectedFiles []string) model.ConnectionTester {
	switch settings.Type {
	case model.FileStorageTypeNone, model.FileStorageTypeDefault:
		return AlwaysHappyConnectionTester{}
	case model.FileStorageTypeLocal:
		return fs.NewFSConnectionTester(settings.Local, expectedFiles)
	case model.FileStorageTypeS3:
		return s3.NewFileStorageConnectionTester(settings.S3, expectedFiles)
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
