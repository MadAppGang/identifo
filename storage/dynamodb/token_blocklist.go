package dynamodb

import (
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/madappgang/identifo/v2/logging"
	"github.com/madappgang/identifo/v2/model"
)

const blacklistedTokensTableName = "BlacklistedTokens"

// NewTokenBlacklist creates new DynamoDB token storage.
func NewTokenBlacklist(
	logger *slog.Logger,
	settings model.DynamoDatabaseSettings,
) (model.TokenBlacklist, error) {
	if len(settings.Endpoint) == 0 || len(settings.Region) == 0 {
		return nil, ErrorEmptyEndpointRegion
	}

	// create database
	db, err := NewDB(settings.Endpoint, settings.Region)
	if err != nil {
		return nil, err
	}

	ts := &TokenBlacklist{
		logger: logger,
		db:     db,
	}
	err = ts.ensureTable()
	return ts, err
}

// TokenBlacklist is a DynamoDB storage for blacklisted tokens.
type TokenBlacklist struct {
	logger *slog.Logger
	db     *DB
}

// ensureTable ensures that token blacklist exists.
func (tb *TokenBlacklist) ensureTable() error {
	exists, err := tb.db.IsTableExists(blacklistedTokensTableName)
	if err != nil {
		return fmt.Errorf("error while checking if %s exists: %w", blacklistedTokensTableName, err)
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
		return fmt.Errorf("error while creating %s table: %w", blacklistedTokensTableName, err)
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
		tb.logger.Error("Error while marshaling token", logging.FieldError, err)
		return ErrorInternalError
	}

	input := &dynamodb.PutItemInput{
		Item:      t,
		TableName: aws.String(blacklistedTokensTableName),
	}

	if _, err = tb.db.C.PutItem(input); err != nil {
		tb.logger.Error("Error while putting token to blacklist", logging.FieldError, err)
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
		tb.logger.Error("Error while fetching token from db", logging.FieldError, err)
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
