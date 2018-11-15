package dynamodb

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/madappgang/identifo/model"
)

const (
	//TokensTableName is collection to store refresh tokens
	TokensTableName = "RefreshTokens"
)

//NewTokenStorage creates mew dynamodb storage
func NewTokenStorage(db *DB) (model.TokenStorage, error) {
	ts := &TokenStorage{db: db}
	err := ts.ensureTable()
	return ts, err

}

//TokenStorage is dynamodb token storage
type TokenStorage struct {
	db *DB
}

//ensureTable ensures app storage table is exists in database
func (ts *TokenStorage) ensureTable() error {
	exists, err := ts.db.isTableExists(TokensTableName)
	if err != nil {
		log.Println(err)
		return err
	}
	if !exists {
		//create table, AWS DynamoDB table creation is overcomplicated for sure
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
			ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(10),
				WriteCapacityUnits: aws.Int64(10),
			},
			TableName: aws.String(TokensTableName),
		}
		_, err = ts.db.C.CreateTable(input)
		fmt.Println("21", err)
		return err
	}
	return nil
}

//SaveToken save token in database
func (ts *TokenStorage) SaveToken(token string) error {
	if len(token) == 0 {
		return model.ErrorWrongDataFormat
	}
	if ts.HasToken(token) {
		return nil
	}

	t := Token{Token: token}
	tv, err := dynamodbattribute.MarshalMap(t)
	if err != nil {
		log.Println(err)
		return ErrorInternalError
	}

	input := &dynamodb.PutItemInput{
		Item:      tv,
		TableName: aws.String(TokensTableName),
	}
	_, err = ts.db.C.PutItem(input)
	if err != nil {
		log.Println(err)
		return ErrorInternalError
	}
	return nil

}

//HasToken returns true if the token in the storage
func (ts *TokenStorage) HasToken(token string) bool {
	if len(token) == 0 {
		return false
	}

	result, err := ts.db.C.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(TokensTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"token": {
				S: aws.String(token),
			},
		},
	})
	if err != nil {
		log.Println(err)
		return false
	}
	//empty result
	if result.Item == nil {
		return false
	}
	return true
}

//RevokeToken removes token from the storage
func (ts *TokenStorage) RevokeToken(token string) error {
	if !ts.HasToken(token) {
		return model.ErrorNotFound
	}
	_, err := ts.db.C.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(TokensTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"token": {
				S: aws.String(token),
			},
		},
	})
	if err != nil {
		log.Println(err)
		return ErrorInternalError
	}
	return nil

}

//Token is struct to store tokens in database
type Token struct {
	Token string `json:"token,omitempty"`
}
