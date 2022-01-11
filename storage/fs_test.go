package storage_test

import (
	"fmt"
	"io/fs"
	"testing"

	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage"
)

func TestNewFSWithLocalFolder(t *testing.T) {
	settings := model.FileStorageSettings{
		Type: model.FileStorageTypeLocal,
		Local: model.FileStorageLocal{
			FolderPath: "../test/artifacts/templates",
		},
	}

	testFSContent(settings, t)
}

func TestNewFSWithS3(t *testing.T) {
	settings := model.FileStorageSettings{
		Type: model.FileStorageTypeS3,
		S3: model.FileStorageS3{
			Region: "ap-southeast-2",
			Bucket: "identifo-public",
			Folder: "test/templates",
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
