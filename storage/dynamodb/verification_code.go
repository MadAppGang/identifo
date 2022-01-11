package dynamodb

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/madappgang/identifo/v2/model"
)

const (
	// verificationCodesTableName is a table name for verification codes.
	verificationCodesTableName = "VerificationCodes"

	// verificationCodesExpirationTime specifies time before deleting records.
	verificationCodesExpirationTime = 5 * time.Minute

	phoneField     = "phone"
	codeField      = "code"
	expiresAtField = "expiresAt"
)

// NewVerificationCodeStorage creates and provisions new DynamoDB verification code storage.
func NewVerificationCodeStorage(settings model.DynamoDatabaseSettings) (model.VerificationCodeStorage, error) {
	if len(settings.Endpoint) == 0 || len(settings.Region) == 0 {
		return nil, ErrorEmptyEndpointRegion
	}

	// create database
	db, err := NewDB(settings.Endpoint, settings.Region)
	if err != nil {
		return nil, err
	}

	vcs := &VerificationCodeStorage{db: db}
	err = vcs.ensureTable()
	return vcs, err
}

// VerificationCodeStorage implements verification code storage interface.
type VerificationCodeStorage struct {
	db *DB
}

// IsVerificationCodeFound checks whether verification code can be found.
func (vcs *VerificationCodeStorage) IsVerificationCodeFound(phone, code string) (bool, error) {
	result, err := vcs.db.C.Query(&dynamodb.QueryInput{
		TableName:              aws.String(verificationCodesTableName),
		KeyConditionExpression: aws.String("phone = :phone"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":phone": {S: aws.String(phone)},
		},
	})
	if err != nil {
		log.Println("Error querying for verification code:", err)
		return false, ErrorInternalError
	}
	if len(result.Items) == 0 {
		return false, nil
	}
	return true, nil
}

// CreateVerificationCode inserts new verification code to the database.
func (vcs *VerificationCodeStorage) CreateVerificationCode(phone, code string) error {
	// Remove old item first.
	delInput := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			phoneField: {S: aws.String(phone)},
		},
		TableName: aws.String(verificationCodesTableName),
	}

	if _, err := vcs.db.C.DeleteItem(delInput); err != nil {
		log.Println("Error deleting old verification code: ", err)
		return ErrorInternalError
	}

	// Then put a new one.
	item, err := dynamodbattribute.MarshalMap(map[string]interface{}{
		phoneField:     phone,
		codeField:      code,
		expiresAtField: time.Now().Add(verificationCodesExpirationTime),
	})
	if err != nil {
		log.Println("Error marshalling verification code:", err)
		return ErrorInternalError
	}

	putInput := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(verificationCodesTableName),
	}

	if _, err := vcs.db.C.PutItem(putInput); err != nil {
		log.Println("Error putting verification code to database:", err)
		return ErrorInternalError
	}
	return err
}

// ensureTable ensures that verification code storage table exists in the database.
func (vcs *VerificationCodeStorage) ensureTable() error {
	exists, err := vcs.db.IsTableExists(verificationCodesTableName)
	if err != nil {
		log.Println("Error checking for verification codes table existence:", err)
		return err
	}
	if exists {
		return nil
	}

	createTableInput := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String(phoneField),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String(phoneField),
				KeyType:       aws.String("HASH"),
			},
		},
		BillingMode: aws.String("PAY_PER_REQUEST"),
		TableName:   aws.String(verificationCodesTableName),
	}

	if _, err = vcs.db.C.CreateTable(createTableInput); err != nil {
		log.Println("Error creating table:", err)
		return err
	}

	ttlSpecification := &dynamodb.TimeToLiveSpecification{
		AttributeName: aws.String(expiresAtField),
		Enabled:       aws.Bool(true),
	}
	ttlInput := &dynamodb.UpdateTimeToLiveInput{
		TableName:               aws.String(verificationCodesTableName),
		TimeToLiveSpecification: ttlSpecification,
	}

	if _, err = vcs.db.C.UpdateTimeToLive(ttlInput); err != nil {
		if err.Error() == dynamodb.ErrCodeResourceNotFoundException {
			// Then Verification Codes table must be in creating status. Let's give it some time.
			for i := 0; i < 5; i++ {
				time.Sleep(5 * time.Second)
				log.Println("Retry setting expiration time...")
				if _, err = vcs.db.C.UpdateTimeToLive(ttlInput); err == nil {
					log.Println("Expiration time successfully set")
					break
				}
			}
		}
	}
	return err
}

// Close does nothing here.
func (vcs *VerificationCodeStorage) Close() {}
