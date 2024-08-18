package dynamodb

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// NewDB creates new database connection.
func NewDB(endpoint string, region string) (*DB, error) {
	config := &aws.Config{
		Region:   aws.String(region),
		Endpoint: aws.String(endpoint),
	}
	sess, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}

	return &DB{C: dynamodb.New(sess)}, nil
}

// DB represents connection to AWS DynamoDB service or local instance.
type DB struct {
	C *dynamodb.DynamoDB
}

// IsTableExists checks if table exists.
func (db *DB) IsTableExists(table string) (bool, error) {
	input := &dynamodb.DescribeTableInput{
		TableName: aws.String(table),
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.C.DescribeTableWithContext(timeoutCtx, input)
	if AwsErrorErrorNotFound(err) {
		return false, nil
		// if table not exists - create table
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (db *DB) DeleteTable(table string) error {
	sess, err := session.NewSession()
	if err != nil {
		return err
	}

	svc := dynamodb.New(sess)
	input := &dynamodb.DeleteTableInput{
		TableName: aws.String(table),
	}

	_, err = svc.DeleteTable(input)
	return err
}

// AwsErrorErrorNotFound checks if error has type dynamodb.ErrCodeResourceNotFoundException.
func AwsErrorErrorNotFound(err error) bool {
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
