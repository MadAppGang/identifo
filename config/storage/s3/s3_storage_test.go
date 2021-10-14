package s3_test

import (
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	s3s "github.com/madappgang/identifo/config/storage/s3"
	"github.com/madappgang/identifo/model"
)

func TestS3ConfigSource(t *testing.T) {
	if os.Getenv("IDENTIFO_TEST_INGRATION") == "" {
		t.SkipNow()
	}

	awsEndpoint := os.Getenv("IDENTIFO_TEST_AWS_ENDPOINT")

	putTestFileTOS3(t, awsEndpoint)

	c, err := s3s.NewConfigurationStorage(model.ConfigStorageSettings{
		Type: model.ConfigStorageTypeS3,
		S3: &model.S3StorageSettings{
			Region:   settings.Region,
			Bucket:   settings.Bucket,
			Key:      settings.Key,
			Endpoint: awsEndpoint,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	settings, err := c.LoadServerSettings(true)
	if err != nil {
		t.Fatal(err)
	}

	if settings.General.Host != "example.com" {
		t.Fatal("wrong config")
	}
}

func putTestFileTOS3(t *testing.T, endpoint string) {
	s3client, err := s3s.NewS3Client(settings.Region, endpoint)
	if err != nil {
		t.Error(err)
		return
	}
	_, _ = s3client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String("identifo-public"),
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
	if err != nil {
		t.Error(err)
	}
}
