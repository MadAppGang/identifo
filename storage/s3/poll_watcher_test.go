package s3_test

import (
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	s3s "github.com/madappgang/identifo/v2/storage/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWatcher(t *testing.T) {
	// localS3Debug()

	if os.Getenv("IDENTIFO_TEST_INTEGRATION") == "" {
		t.SkipNow()
	}

	settings.Endpoint = os.Getenv("IDENTIFO_TEST_AWS_ENDPOINT")

	s3client, err := s3s.NewS3Client(settings.Region, settings.Endpoint)
	require.NoError(t, err)

	_, _ = s3client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(settings.Bucket),
	})

	watcher, err := s3s.NewPollWatcher(settings, time.Second*2)
	require.NoError(t, err)

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
	// localS3Debug()
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
		Bucket: aws.String(settings.Bucket),
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
