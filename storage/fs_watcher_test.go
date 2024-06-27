package storage_test

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
	"github.com/madappgang/identifo/v2/storage"
	s3s "github.com/madappgang/identifo/v2/storage/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const s3Region = "ap-southeast-2"

func TestNewFSWatcher(t *testing.T) {
	settings := model.FileStorageSettings{
		Type: model.FileStorageTypeLocal,
		Local: model.FileStorageLocal{
			Path: "../fs_test",
		},
	}
	fss, err := storage.NewFS(settings)
	require.NoError(t, err)
	makeFiles()
	defer deleteFiles()

	fileChanged := false
	updatedCount := 0
	changedFiles := []string{}

	time.Sleep(time.Second * 1)

	watcher := storage.NewFSWatcher(fss, []string{"test1.txt", "test2.txt"}, time.Second*2)
	watcher.Watch()

	go func() {
		time.Sleep(time.Second * 1)
		go updateFile("../fs_test/test1.txt")
		go updateFile("../fs_test/test2.txt")
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
	assert.Contains(t, changedFiles, "test2.txt")
	assert.NotContains(t, changedFiles, "file2.txt")
}

func TestNewFSWatcherS3(t *testing.T) {
	// forceS3()
	ep := os.Getenv("IDENTIFO_TEST_AWS_ENDPOINT")
	if ep == "" {
		t.SkipNow()
	}
	settings := model.FileStorageSettings{
		Type: model.FileStorageTypeS3,
		S3: model.FileStorageS3{
			Region:   "ap-southeast-2",
			Bucket:   "identifo",
			Key:      "/fs_test",
			Endpoint: ep,
		},
	}
	fss, err := storage.NewFS(settings)
	require.NoError(t, err)

	s3client := getS3Client(t, ep)
	fileSettings := settings.S3
	fileSettings.Key = "/fs_test/watched.txt"
	uploadS3File(t, s3client, fileSettings, fileSettings.Key)

	fileChanged := false
	updatedCount := 0
	changedFiles := []string{}

	time.Sleep(time.Second * 1)

	watcher := storage.NewFSWatcher(fss, []string{"watched.txt", "watched_new.txt"}, time.Second*2)
	watcher.Watch()

	go func() {
		time.Sleep(time.Second * 1)
		go uploadS3File(t, s3client, fileSettings, "/fs_test/watched_new.txt")
		go uploadS3File(t, s3client, fileSettings, "/fs_test/watched.txt")
		go uploadS3File(t, s3client, fileSettings, "/fs_test/unwatched.txt")
	}()

outFor:
	for {
		select {
		// case err := <-watcher.ErrorChan():
		// 	assert.NoError(t, err)
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
	assert.Contains(t, changedFiles, "watched.txt")
	assert.Contains(t, changedFiles, "watched_new.txt")
	assert.NotContains(t, changedFiles, "unwatched.txt")
}

func makeFiles() {
	os.Mkdir("../fs_test", os.ModePerm)
	f, _ := os.Create("../fs_test/test1.txt")
	f.WriteString("Hello")
	f.Close()
	f, _ = os.Create("../fs_test/test2.txt")
	f.WriteString("Hello2")
	f.Close()
	f, _ = os.Create("../fs_test/test3.txt")
	f.WriteString("Hello3")
	f.Close()
}

func updateFile(name string) {
	data := fmt.Sprintf("This is file data has been created at: %v", time.Now())
	_ = os.WriteFile(name, []byte(data), 0o644)
}

func deleteFiles() {
	os.RemoveAll("../fs_test")
}

func getS3Client(t *testing.T, endpoint string) *s3.S3 {
	s3client, err := s3s.NewS3Client(s3Region, endpoint)
	require.NoError(t, err)
	_, err = s3client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(s3Region),
	})
	if err != nil && !strings.Contains(err.Error(), "BucketAlreadyOwnedByYou: Your previous request to create the named bucket succeeded and you already own it.") {
		require.NoError(t, err)
	}
	return s3client
}

func uploadS3File(t *testing.T, s3client *s3.S3, s model.FileStorageS3, key string) {
	newFilecontent := []byte(fmt.Sprintf("This content has been changed at %v", time.Now().Unix()))
	input := &s3.PutObjectInput{
		Bucket:             aws.String("identifo"),
		Key:                aws.String(key),
		Body:               bytes.NewReader(newFilecontent),
		ContentDisposition: aws.String("attachment"),
	}
	_, err := s3client.PutObject(input)
	assert.NoError(t, err)
}

func forceS3() {
	os.Setenv("IDENTIFO_TEST_INTEGRATION", "1")
	os.Setenv("IDENTIFO_TEST_AWS_ENDPOINT", "http://localhost:9000")
	os.Setenv("AWS_ACCESS_KEY_ID", "testing")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "testing_secret")
	os.Setenv("IDENTIFO_FORCE_S3_PATH_STYLE", "1")
}

func TestKeysWithFixedSlashed(t *testing.T) {
	data := storage.KeysWithFixedSlashed([]string{
		"./folder/path.txt",
		"/folder/path.txt",
		"/folder/path.txt/",
		"..//folder/path.txt",
		"///////////",
	})
	assert.Len(t, data, 5)
	assert.Equal(t, "./folder/path.txt", data[0])
	assert.Equal(t, "folder/path.txt", data[1])
	assert.Equal(t, "folder/path.txt", data[2])
	assert.Equal(t, "..//folder/path.txt", data[3])
	assert.Equal(t, "/////////", data[4])
}
