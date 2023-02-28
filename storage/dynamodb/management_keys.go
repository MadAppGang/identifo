package dynamodb

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/madappgang/identifo/v2/model"
)

const (
	managementKeyTableName = "ManagementKeys"
)

// NewManagementKeysStorage creates and provisions new management keys storage instance.
func NewManagementKeysStorage(settings model.DynamoDatabaseSettings) (model.ManagementKeysStorage, error) {
	if len(settings.Endpoint) == 0 || len(settings.Region) == 0 {
		return nil, ErrorEmptyEndpointRegion
	}

	// create database
	db, err := NewDB(settings.Endpoint, settings.Region)
	if err != nil {
		return nil, err
	}

	us := &ManagementKeysStorage{db: db}
	err = us.ensureTable()
	return us, err
}

// UserStorage stores and manages data in DynamoDB storage.
type ManagementKeysStorage struct {
	db *DB
}

func (ms *ManagementKeysStorage) GetKey(ctx context.Context, id string) (model.ManagementKey, error) {
	return model.ManagementKey{}, errors.New("not implemented")
}

func (ms *ManagementKeysStorage) CreateKey(ctx context.Context, name string, scopes []string) (model.ManagementKey, error) {
	return model.ManagementKey{}, errors.New("not implemented")
}

func (ms *ManagementKeysStorage) DisableKey(ctx context.Context, id string) (model.ManagementKey, error) {
	return model.ManagementKey{}, errors.New("not implemented")
}

func (ms *ManagementKeysStorage) RenameKey(ctx context.Context, id, name string) (model.ManagementKey, error) {
	return model.ManagementKey{}, errors.New("not implemented")
}

func (ms *ManagementKeysStorage) ChangeScopesForKey(ctx context.Context, id string, scopes []string) (model.ManagementKey, error) {
	return model.ManagementKey{}, errors.New("not implemented")
}

func (ms *ManagementKeysStorage) UseKey(ctx context.Context, id string) (model.ManagementKey, error) {
	return model.ManagementKey{}, errors.New("not implemented")
}

func (ms *ManagementKeysStorage) GeyAllKeys(ctx context.Context) ([]model.ManagementKey, error) {
	return nil, errors.New("not implemented")
}

// ensureTable ensures that user storage table exists in the database.
// I'm hiding it in the end of the file, because AWS devs, you are killing me with this API.
func (ms *ManagementKeysStorage) ensureTable() error {
	exists, err := ms.db.IsTableExists(managementKeyTableName)
	if err != nil {
		log.Println("Error checking for table existence:", err)
		return err
	}
	if !exists {
		// create table, AWS DynamoDB table creation is overcomplicated for sure
		input := &dynamodb.CreateTableInput{
			AttributeDefinitions: []*dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String("id"),
					AttributeType: aws.String("S"),
				},
				{
					AttributeName: aws.String("name"),
					AttributeType: aws.String("S"),
				},
				{
					AttributeName: aws.String("active"),
					AttributeType: aws.String("BOOL"),
				},
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("id"),
					KeyType:       aws.String("HASH"),
				},
			},

			BillingMode: aws.String("PAY_PER_REQUEST"),
			TableName:   aws.String(managementKeyTableName),
		}
		if _, err = ms.db.C.CreateTable(input); err != nil {
			log.Println("Error creating table:", err)
			return err
		}
	}

	return nil
}
