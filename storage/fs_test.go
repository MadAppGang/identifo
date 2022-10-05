package storage_test

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage"
	"github.com/stretchr/testify/require"
)

const (
	templPath  = "../test/artifacts/templates"
	testRegion = "ap-southeast-2"
	testBucket = "identifo-public"
	testFolder = "test/templates"
)

func TestNewFSWithLocalFolder(t *testing.T) {
	settings := model.FileStorageSettings{
		Type: model.FileStorageTypeLocal,
		Local: model.FileStorageLocal{
			FolderPath: templPath,
		},
	}

	testFSContent(settings, t)
}

func TestNewFSWithS3(t *testing.T) {
	ep := os.Getenv("IDENTIFO_TEST_AWS_ENDPOINT")
	if ep != "" {
		// this is for local tests
		putTestFileTOS3(t, ep)
	}

	settings := model.FileStorageSettings{
		Type: model.FileStorageTypeS3,
		S3: model.FileStorageS3{
			Region:   testRegion,
			Bucket:   testBucket,
			Folder:   testFolder,
			Endpoint: ep,
		},
	}

	testFSContent(settings, t)
}

func testFSContent(settings model.FileStorageSettings, t *testing.T) {
	fss, err := storage.NewFS(settings)
	if err != nil {
		t.Fatalf("error creating local fs with error: %v", err)
	}

	// print out all files in root.
	printFolderContent(fss, ".")

	file, err := fss.Open("mail1.template")
	if err != nil {
		t.Fatalf("error opening email template 1: %v", err)
	}

	stat, err := file.Stat()
	if err != nil {
		t.Fatalf("error opening getting stat for email template 1: %v", err)
	}

	if stat.Size() <= 0 {
		if err != nil {
			t.Fatalf("email template is empty 1: %v", err)
		}
	}

	fdata, err := fs.ReadFile(fss, "mail1.template")
	if err != nil {
		t.Fatalf("error reading email template 1: %v", err)
	}
	fmt.Printf("data: %s\n", string(fdata))
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

func putTestFileTOS3(t *testing.T, endpoint string) {
	cfg := aws.NewConfig().
		WithEndpoint(endpoint).
		WithRegion(testRegion)

	sess, err := session.NewSession(cfg)
	require.NoError(t, err)

	s := s3.New(sess, cfg)

	_, _ = s.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String("identifo-public"),
	})

	d, err := os.ReadFile(filepath.Join(templPath, "mail1.template"))
	require.NoError(t, err)

	input := &s3.PutObjectInput{
		Bucket:             aws.String(testBucket),
		Key:                aws.String(filepath.Join(testFolder, "mail1.template")),
		Body:               bytes.NewReader(d),
		ContentDisposition: aws.String("attachment"),
	}
	_, err = s.PutObject(input)
	require.NoError(t, err)
}
