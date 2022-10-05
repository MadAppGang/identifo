package s3_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	s3s "github.com/madappgang/identifo/v2/config/storage/s3"

	"github.com/madappgang/identifo/v2/model"
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

var settings = model.S3StorageSettings{
	Region: "ap-southeast-2",
	Bucket: "identifo-public",
	Key:    "test/config-boltdb.yaml",
}

func TestWatcher(t *testing.T) {
	if os.Getenv("IDENTIFO_TEST_INTEGRATION") == "" {
		t.SkipNow()
	}

	settings.Endpoint = os.Getenv("IDENTIFO_TEST_AWS_ENDPOINT")

	s3client, err := s3s.NewS3Client(settings.Region, settings.Endpoint)
	if err != nil {
		t.Error(err)
		return
	}
	_, _ = s3client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String("identifo-public"),
	})

	watcher, err := s3s.NewPollWatcher(settings, time.Second*4)
	if err != nil {
		t.Fatalf("unable to create watcher with error: %v", err)
	}

	fileChanged := false
	watcher.Watch()

	go func() {
		time.Sleep(time.Second * 1)

		newFilecontent := []byte(fmt.Sprintf("This content has been changed at %v", time.Now().Unix()))
		input := &s3.PutObjectInput{
			Bucket:             aws.String(settings.Bucket),
			Key:                aws.String(settings.Key),
			Body:               bytes.NewReader(newFilecontent),
			ContentDisposition: aws.String("attachment"),
		}
		_, err = s3client.PutObject(input)
		if err != nil {
			t.Error(err)
		}
	}()

	select {
	case err := <-watcher.ErrorChan():
		t.Error(err)
		return
	case <-watcher.WatchChan():
		fileChanged = true
		t.Log("getting file changed")
	case <-time.After(time.Second * 30):
		t.Error("timeout waiting file update")
	}

	// let's check if file has been changed after select finished
	if !fileChanged {
		t.Error("no file change handled")
	}
}
