package dynamodb

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (

	//DefaultRegion default AWS region, full list is here: https://docs.aws.amazon.com/general/latest/gr/rande.html
	DefaultRegion = "us-west-2"
)

//NewDB creates new database connection
func NewDB(endpoint string, region string) (*DB, error) {
	r := DefaultRegion
	if len(region) > 0 {
		r = region
	}
	config := &aws.Config{
		Region:   aws.String(r),
		Endpoint: aws.String(endpoint),
	}
	sess, err := session.NewSession(config)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	m := DB{}
	m.C = dynamodb.New(sess)
	return &m, nil
}

//DB represents connection to AWS DynamoDB service or local instance
type DB struct {
	C *dynamodb.DynamoDB
}

//isTableExists checks if table exists database
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

//awsErrorErrorNotFound check general error type to be dynamodb.ErrCodeResourceNotFoundException
func awsErrorErrorNotFound(err error) bool {
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == dynamodb.ErrCodeResourceNotFoundException {
				return true
			}
		}
	}
	return false
}
