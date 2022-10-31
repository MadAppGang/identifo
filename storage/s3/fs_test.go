package s3_test

import (
	"fmt"
	"io/fs"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testFolder = "test/templates"
)

func TestNewFSWithS3(t *testing.T) {
	// localS3Debug()
	ep := os.Getenv("IDENTIFO_TEST_AWS_ENDPOINT")
	if ep == "" {
		t.SkipNow()
	}

	s3client := getS3Client(t, ep)
	uploadS3File(t, s3client, settings, "/test/templates/file_whatever.txt")

	settings := model.FileStorageSettings{
		Type: model.FileStorageTypeS3,
		S3: model.FileStorageS3{
			Region:   settings.Region,
			Bucket:   settings.Bucket,
			Key:      testFolder,
			Endpoint: ep,
		},
	}

	testFSContent(settings, s3client, t)
}

func testFSContent(sts model.FileStorageSettings, s3client *s3.S3, t *testing.T) {
	fss, err := storage.NewFS(sts)
	if err != nil {
		t.Fatalf("error creating local fs with error: %v", err)
	}

	// print out all files in root.
	printFolderContent(fss, ".")

	file, err := fss.Open("file_whatever.txt")
	assert.NoError(t, err)

	stat, err := file.Stat()
	assert.NoError(t, err)

	assert.Greater(t, stat.Size(), int64(0))

	fdata, err := fs.ReadFile(fss, "file_whatever.txt")
	assert.NoError(t, err)
	assert.NotEmpty(t, fdata)
}

func printFolderContent(fss fs.FS, path string) {
	_ = fs.WalkDir(fss, path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			fmt.Println("dir:", path)
			return nil
		}
		fmt.Println("file:", path)
		return nil
	})
}

func TestNewS3FSPath(t *testing.T) {
	v := fs.ValidPath("/server-config-prod.yaml")
	require.False(t, v)

	v2 := fs.ValidPath("server-config-prod.yaml")
	require.True(t, v2)
}
