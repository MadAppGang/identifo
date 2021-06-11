package storage

import (
	"fmt"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/storage/dynamodb"
	"github.com/madappgang/identifo/storage/fs"
	"github.com/madappgang/identifo/storage/s3"
)

var defaultStaticFileFolder = "./static"

func NewStaticFileStorage(settings model.StaticFilesStorageSettings) (model.StaticFilesStorage, error) {
	fallback, err := newDefaultStaticFileStorage()
	if err != nil {
		return nil, err
	}
	switch settings.Type {
	case model.StaticFilesStorageTypeLocal:
		return fs.NewStaticFilesStorage(settings.Local, fallback)
	case model.StaticFilesStorageTypeDynamoDB:
		return dynamodb.NewStaticFilesStorage(settings.Dynamo, fallback)
	case model.StaticFilesStorageTypeS3:
		return s3.NewStaticFilesStorage(settings.S3, fallback)
	default:
		// always return a fallback
		fmt.Printf("unable to create static storage with type %s, using fallback readonly folder", settings.Type)
		return fallback, nil
	}
}

// newDefaultStaticFileStorage creates default file storage
func newDefaultStaticFileStorage() (model.StaticFilesStorage, error) {
	return fs.DefaultFileStorage(defaultStaticFileFolder)
}
