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
	"github.com/stretchr/testify/assert"

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
	// localDebug()

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

	watcher, err := s3s.NewPollWatcher(settings, time.Second*2)
	if err != nil {
		t.Fatalf("unable to create watcher with error: %v", err)
	}
	uploadS3File(t, s3client, settings, settings.Key)
	// wait local S3 to upload the file
	time.Sleep(time.Second * 1)

	fileChanged := false
	watcher.Watch()

	go func() {
		time.Sleep(time.Second * 1)
		go uploadS3File(t, s3client, settings, settings.Key)
	}()

	updatedCount := 0
outFor:
	for {
		select {
		case err := <-watcher.ErrorChan():
			t.Error(err)
			return
		case <-watcher.WatchChan():
			updatedCount++
			fileChanged = true
			t.Log("getting file changed")
		case <-time.After(time.Second * 10):
			// wait for all updates here for 5 secs
			break outFor
		}
	}

	assert.True(t, fileChanged)
	// should fire update only once!
	assert.Equal(t, updatedCount, 1)
}

func TestWatcherWithListOfFiles(t *testing.T) {
	// localDebug()
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

	watcher, err := s3s.NewPollWatcherForKeyList(settings, []string{"file1.txt", "file2.txt", "file3.txt"}, time.Second*2)
	if err != nil {
		t.Fatalf("unable to create watcher with error: %v", err)
	}
	uploadS3File(t, s3client, settings, "file1.txt")
	uploadS3File(t, s3client, settings, "file2.txt")
	uploadS3File(t, s3client, settings, "file3.txt")
	// wait local S3 to upload the files
	time.Sleep(time.Second * 1)

	fileChanged := false
	watcher.Watch()

	go func() {
		time.Sleep(time.Second * 1)
		go uploadS3File(t, s3client, settings, "file1.txt")
		go uploadS3File(t, s3client, settings, "file3.txt")
	}()

	updatedCount := 0
	changedFiles := []string{}
outFor:
	for {
		select {
		case err := <-watcher.ErrorChan():
			t.Error(err)
			return
		case files := <-watcher.WatchChan():
			fileChanged = true
			updatedCount++
			changedFiles = append(changedFiles, files...)
			t.Logf("getting file changed: %+v", files)
		case <-time.After(time.Second * 5):
			// wait and collect data for 5 seconds and then breaks loop
			break outFor
		}
	}

	assert.True(t, fileChanged)
	// should fire update only once!
	assert.LessOrEqual(t, updatedCount, 2) // could update in one butch or two bunches
	assert.GreaterOrEqual(t, updatedCount, 1)
	assert.Contains(t, changedFiles, "file1.txt")
	assert.Contains(t, changedFiles, "file3.txt")
	assert.NotContains(t, changedFiles, "file2.txt")
}

func uploadS3File(t *testing.T, s3client *s3.S3, s model.S3StorageSettings, key string) {
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

func localDebug() {
	os.Setenv("IDENTIFO_TEST_INTEGRATION", "1")
	os.Setenv("IDENTIFO_TEST_AWS_ENDPOINT", "http://localhost:5001")
	os.Setenv("AWS_ACCESS_KEY_ID", "testing")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "testing")
	os.Setenv("AWS_SECURITY_TOKEN", "testing")
	os.Setenv("AWS_SESSION_TOKEN", "testing")
}
