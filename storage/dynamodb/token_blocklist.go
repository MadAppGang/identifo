package dynamodb

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/madappgang/identifo/v2/model"
)

const blacklistedTokensTableName = "BlacklistedTokens"

// NewTokenBlacklist creates new DynamoDB token storage.
func NewTokenBlacklist(settings model.DynamoDatabaseSettings) (model.TokenBlacklist, error) {
	if len(settings.Endpoint) == 0 || len(settings.Region) == 0 {
		return nil, ErrorEmptyEndpointRegion
	}

	// create database
	db, err := NewDB(settings.Endpoint, settings.Region)
	if err != nil {
		return nil, err
	}

	ts := &TokenBlacklist{db: db}
	err = ts.ensureTable()
	return ts, err
}

// TokenBlacklist is a DynamoDB storage for blacklisted tokens.
type TokenBlacklist struct {
	db *DB
}

// ensureTable ensures that token blacklist exists.
func (tb *TokenBlacklist) ensureTable() error {
	exists, err := tb.db.IsTableExists(blacklistedTokensTableName)
	if err != nil {
		log.Printf("Error while checking if %s exists: %v", blacklistedTokensTableName, err)
		return err
	}
	if exists {
		return nil
	}

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("token"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("token"),
				KeyType:       aws.String("HASH"),
			},
		},
		BillingMode: aws.String("PAY_PER_REQUEST"),
		TableName:   aws.String(blacklistedTokensTableName),
	}

	if _, err = tb.db.C.CreateTable(input); err != nil {
		log.Printf("Error while creating %s table: %v", blacklistedTokensTableName, err)
		return err
	}
	return nil
}

// Add adds token to the blacklist.
func (tb *TokenBlacklist) Add(token string) error {
	if len(token) == 0 {
		return model.ErrorWrongDataFormat
	}
	if tb.IsBlacklisted(token) {
		return nil
	}

	t, err := dynamodbattribute.MarshalMap(Token{Token: token})
	if err != nil {
		log.Println(err)
		return ErrorInternalError
	}

	input := &dynamodb.PutItemInput{
		Item:      t,
		TableName: aws.String(blacklistedTokensTableName),
	}

	if _, err = tb.db.C.PutItem(input); err != nil {
		log.Println("Error while putting token to blacklist:", err)
		return ErrorInternalError
	}
	return nil
}

// IsBlacklisted returns true if token is blacklisted.
func (tb *TokenBlacklist) IsBlacklisted(token string) bool {
	if len(token) == 0 {
		return false
	}

	result, err := tb.db.C.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(blacklistedTokensTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"token": {
				S: aws.String(token),
			},
		},
	})
	if err != nil {
		log.Println("Error while fetching token from db:", err)
		return false
	}

	// Empty result.
	if result.Item == nil {
		return false
	}
	return true
}

// Close does nothing here.
func (tb *TokenBlacklist) Close() {}
