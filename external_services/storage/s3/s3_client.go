package s3

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	awsAccessKeyEnvName       = "AWS_ACCESS_KEY"
	awsSecretAccessKeyEnvName = "AWS_SECRET_ACCESS_KEY"
)

// NewS3Client creates and returns new S3 client.
func NewS3Client(region string) (*s3.S3, error) {
	awsAccessKey := os.Getenv(awsAccessKeyEnvName)
	if len(awsAccessKey) == 0 {
		return nil, fmt.Errorf("No %s specified", awsAccessKeyEnvName)
	}
	awsSecret := os.Getenv(awsSecretAccessKeyEnvName)
	if len(awsSecretAccessKeyEnvName) == 0 {
		return nil, fmt.Errorf("No %s specified", awsSecretAccessKeyEnvName)
	}

	if len(region) == 0 {
		return nil, fmt.Errorf("No S3 region for configuration specified")
	}

	creds := credentials.NewStaticCredentials(awsAccessKey, awsSecret, "")
	cfg := aws.NewConfig().WithRegion(region).WithCredentials(creds)

	sess, err := session.NewSession(cfg)
	if err != nil {
		return nil, fmt.Errorf("Cannot create new session: %s", err)
	}
	return s3.New(sess, cfg), nil
}
