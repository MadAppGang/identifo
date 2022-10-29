package s3_test

import (
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/madappgang/identifo/v2/model"
	s3s "github.com/madappgang/identifo/v2/storage/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestS3ConfigSource(t *testing.T) {
	if os.Getenv("IDENTIFO_TEST_INTEGRATION") == "" {
		t.SkipNow()
	}

	awsEndpoint := os.Getenv("IDENTIFO_TEST_AWS_ENDPOINT")

	putTestFileTOS3(t, awsEndpoint)

	c, err := s3s.NewConfigurationStorage(model.FileStorageSettings{
		Type: model.FileStorageTypeS3,
		S3: model.FileStorageS3{
			Region:   settings.Region,
			Bucket:   settings.Bucket,
			Key:      settings.Key,
			Endpoint: awsEndpoint,
		},
	})
	require.NoError(t, err)
	settings, errs := c.LoadServerSettings(true)
	require.Empty(t, errs)
	assert.Equal(t, settings.General.Host, "example.com")
}

func putTestFileTOS3(t *testing.T, endpoint string) {
	s3client, err := s3s.NewS3Client(settings.Region, endpoint)
	require.NoError(t, err)

	_, _ = s3client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(settings.Bucket),
	})

	configFile := `
general:
    host: example.com`

	input := &s3.PutObjectInput{
		Bucket:             aws.String(settings.Bucket),
		Key:                aws.String(settings.Key),
		Body:               strings.NewReader(configFile),
		ContentDisposition: aws.String("attachment"),
	}

	_, err = s3client.PutObject(input)
	require.NoError(t, err)
}
