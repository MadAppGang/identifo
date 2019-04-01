package dynamodb

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// NewDB creates new database connection.
func NewDB(endpoint string, region string) (*DB, error) {
	if len(region) == 0 {
		return nil, errors.New("Empty region string")
	}
	config := &aws.Config{
		Region:   aws.String(region),
		Endpoint: aws.String(endpoint),
	}
	sess, err := session.NewSession(config)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &DB{C: dynamodb.New(sess)}, nil
}

// DB represents connection to AWS DynamoDB service or local instance.
type DB struct {
	C *dynamodb.DynamoDB
}

// isTableExists checks if table exists.
func (db *DB) isTableExists(table string) (bool, error) {
	input := &dynamodb.DescribeTableInput{
		TableName: aws.String(table),
	}
	_, err := db.C.DescribeTable(input)
	if awsErrorErrorNotFound(err) {
		return false, nil
		//if table not exists - create table
	}
	if err != nil {
		log.Println(err)
		return false, err
	}

	return true, nil
}

// awsErrorErrorNotFound checks if error has type dynamodb.ErrCodeResourceNotFoundException.
func awsErrorErrorNotFound(err error) bool {
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
