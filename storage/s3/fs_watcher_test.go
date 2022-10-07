package s3_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFSWatcher(t *testing.T) {
	// localS3Debug()
	ep := os.Getenv("IDENTIFO_TEST_AWS_ENDPOINT")
	if ep == "" {
		t.SkipNow()
	}

	s := model.FileStorageSettings{
		Type: model.FileStorageTypeS3,
		S3: model.FileStorageS3{
			Region:   settings.Region,
			Bucket:   settings.Bucket,
			Key:      testFolder,
			Endpoint: ep,
		},
	}
	fss, err := storage.NewFS(s)
	require.NoError(t, err)
	s3client := getS3Client(t, ep)
	makeFiles(t, s3client)

	fileChanged := false
	updatedCount := 0
	changedFiles := []string{}

	time.Sleep(time.Second * 1)

	watcher := storage.NewFSWatcher(fss, []string{"test1.txt", "test2.txt", "test3.txt"}, time.Second*2)
	watcher.Watch()

	go func() {
		time.Sleep(time.Second * 1)
		go uploadS3File(t, s3client, settings, filepath.Join(testFolder, "test1.txt"))
		go uploadS3File(t, s3client, settings, filepath.Join(testFolder, "test3.txt"))
	}()

outFor:
	for {
		select {
		case err := <-watcher.ErrorChan():
			t.Error(err)
			return
		case files := <-watcher.WatchChan():
			updatedCount++
			fileChanged = true
			changedFiles = append(changedFiles, files...)
			t.Log("getting file changed")
		case <-time.After(time.Second * 5):
			// wait for all updates here for 5 secs
			break outFor
		}
	}

	assert.True(t, fileChanged)
	// should fire update only once!
	assert.LessOrEqual(t, updatedCount, 2) // could update in one butch or two bunches
	assert.GreaterOrEqual(t, updatedCount, 1)
	assert.Contains(t, changedFiles, "test1.txt")
	assert.Contains(t, changedFiles, "test3.txt")
	assert.NotContains(t, changedFiles, "test2.txt")
}

func makeFiles(t *testing.T, s3client *s3.S3) {
	uploadS3File(t, s3client, settings, filepath.Join(testFolder, "test1.txt"))
	uploadS3File(t, s3client, settings, filepath.Join(testFolder, "test2.txt"))
	uploadS3File(t, s3client, settings, filepath.Join(testFolder, "test3.txt"))
}
