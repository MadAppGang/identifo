package dynamodb

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/madappgang/identifo/v2/model"
)

const (
	adminSessionsTableName = "AdminSessions"
)

// DynamoDBSessionStorage is a DynamoDB-backed storage for admin sessions.
type DynamoDBSessionStorage struct {
	db *dynamodb.DynamoDB
}

// NewSessionStorage creates new DynamoDB session storage.
func NewSessionStorage(settings model.DynamoDatabaseSettings) (model.SessionStorage, error) {
	config := &aws.Config{
		Region:   aws.String(settings.Region),
		Endpoint: aws.String(settings.Endpoint),
	}

	sess, err := session.NewSession(config)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	dss := &DynamoDBSessionStorage{db: dynamodb.New(sess)}
	err = dss.ensureTable()
	return dss, err
}

// GetSession fetches session by ID.
func (dss *DynamoDBSessionStorage) GetSession(id string) (model.Session, error) {
	var session model.Session
	if len(id) == 0 {
		return session, model.ErrSessionNotFound
	}

	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String(adminSessionsTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}

	result, err := dss.db.GetItem(getItemInput)
	if err != nil {
		return session, err
	}
	if result.Item == nil {
		return session, model.ErrSessionNotFound
	}

	if err = dynamodbattribute.UnmarshalMap(result.Item, &session); err != nil {
		return session, fmt.Errorf("Error unmarshaling item: %s", err)
	}
	return session, nil
}

// InsertSession inserts session to the storage.
func (dss *DynamoDBSessionStorage) InsertSession(session model.Session) error {
	sessMap, err := dynamodbattribute.MarshalMap(session)
	if err != nil {
		return fmt.Errorf("Error marshalling session: %s", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(adminSessionsTableName),
		Item:      sessMap,
	}

	if _, err = dss.db.PutItem(input); err != nil {
		return fmt.Errorf("Error putting session to the storage: %s", err)
	}
	return err
}

// DeleteSession deletes session from the storage.
func (dss *DynamoDBSessionStorage) DeleteSession(id string) error {
	if len(id) == 0 {
		return nil
	}

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(adminSessionsTableName),
	}
	_, err := dss.db.DeleteItem(input)
	return err
}

// ProlongSession sets new duration for the existing session.
func (dss *DynamoDBSessionStorage) ProlongSession(id string, newDuration model.SessionDuration) error {
	session, err := dss.GetSession(id)
	if err != nil {
		return err
	}

	session.ExpirationTime = time.Now().Add(newDuration.Duration).Unix()

	return dss.InsertSession(session)
}

// ensureTable ensures that admin sessions table exists in database.
func (dss *DynamoDBSessionStorage) ensureTable() error {
	exists, err := dss.isTableExists(adminSessionsTableName)
	if err != nil {
		log.Println("Error checking admins sessions table existence:", err)
		return err
	}
	if exists {
		return nil
	}

	input := &dynamodb.CreateTableInput{
		TableName: aws.String(adminSessionsTableName),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("HASH"),
			},
		},
		BillingMode: aws.String("PAY_PER_REQUEST"),
	}

	_, err = dss.db.CreateTable(input)
	if err != nil {
		return fmt.Errorf("Cannot create admin sessions table: %s", err)
	}

	ttlSpecification := &dynamodb.TimeToLiveSpecification{
		AttributeName: aws.String("expiration_time"),
		Enabled:       aws.Bool(true),
	}
	ttlInput := &dynamodb.UpdateTimeToLiveInput{
		TableName:               aws.String(adminSessionsTableName),
		TimeToLiveSpecification: ttlSpecification,
	}

	if _, err = dss.db.UpdateTimeToLive(ttlInput); err != nil {
		if err.Error() == dynamodb.ErrCodeResourceNotFoundException {
			// Then table must be in creating status. Let's give it some time.
			for i := 0; i < 5; i++ {
				time.Sleep(5 * time.Second)
				log.Println("Retry setting expiration time...")
				if _, err = dss.db.UpdateTimeToLive(ttlInput); err == nil {
					log.Println("Expiration time successfully set")
					break
				}
			}
		}
	}

	return err
}

// isTableExists checks if table exists.
func (dss *DynamoDBSessionStorage) isTableExists(table string) (bool, error) {
	input := &dynamodb.DescribeTableInput{
		TableName: aws.String(table),
	}

	_, err := dss.db.DescribeTable(input)
	if awsErrorNotFound(err) {
		return false, nil
	}
	if err != nil {
		log.Println(err)
		return false, err
	}

	return true, nil
}

// awsErrorNotFound checks if error has type dynamodb.ErrCodeResourceNotFoundException.
func awsErrorNotFound(err error) bool {
	if err == nil {
		return false
	}
	if aerr, ok := err.(awserr.Error); ok {
		if aerr.Code() == dynamodb.ErrCodeResourceNotFoundException {
			return true
		}
	}
	return false
}
func (is *DynamoDBSessionStorage) Close() {

}
