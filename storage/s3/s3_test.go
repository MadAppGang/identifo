package s3_test

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/madappgang/identifo/v2/model"
	s3s "github.com/madappgang/identifo/v2/storage/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// to start the test you need to configure AWS Key and Secret
//
// When you initialize a new service client without providing any credential arguments,
// the SDK uses the default credential provider chain to find AWS credentials.
// The SDK uses the first provider in the chain that returns credentials without an error.
// The default provider chain looks for credentials in the following order:
//
// 1. Environment variables.
// 2. Shared credentials file.
// 3. If your application uses an ECS task definition or RunTask API operation, IAM role for tasks.
// 4. If your application is running on an Amazon EC2 instance, IAM role for Amazon EC2.

var settings = model.FileStorageS3{
	Region: "ap-southeast-2",
	Bucket: "identifo",
	Key:    "test/config-boltdb.yaml",
}

func getS3Client(t *testing.T, endpoint string) *s3.S3 {
	s3client, err := s3s.NewS3Client(settings.Region, endpoint)
	require.NoError(t, err)
	_, err = s3client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(settings.Bucket),
	})
	if err != nil && !strings.Contains(err.Error(), "BucketAlreadyOwnedByYou: Your previous request to create the named bucket succeeded and you already own it.") {
		require.NoError(t, err)
	}
	return s3client
}

func uploadS3File(t *testing.T, s3client *s3.S3, _ model.FileStorageS3, key string) {
	newFilecontent := []byte(fmt.Sprintf("This content has been changed at %v", time.Now().Unix()))
	input := &s3.PutObjectInput{
		Bucket:             aws.String(settings.Bucket),
		Key:                aws.String(key),
		Body:               bytes.NewReader(newFilecontent),
		ContentDisposition: aws.String("attachment"),
	}
	_, err := s3client.PutObject(input)
	assert.NoError(t, err)
}

func localS3Debug() {
	os.Setenv("IDENTIFO_TEST_INTEGRATION", "1")
	os.Setenv("IDENTIFO_TEST_AWS_ENDPOINT", "http://localhost:9000")
	os.Setenv("AWS_ACCESS_KEY_ID", "testing")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "testing_secret")
	os.Setenv("IDENTIFO_FORCE_S3_PATH_STYLE", "1")
}
