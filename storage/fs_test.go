package storage_test

import (
	"fmt"
	"io/fs"
	"testing"

	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	templPath = "../test/artifacts/templates"
)

func TestNewFSWithLocalFolder(t *testing.T) {
	settings := model.FileStorageSettings{
		Type: model.FileStorageTypeLocal,
		Local: model.FileStorageLocal{
			Path: templPath,
		},
	}

	testFSContent(settings, t)
}

func testFSContent(settings model.FileStorageSettings, t *testing.T) {
	fss, err := storage.NewFS(settings)
	require.NoError(t, err)

	// print out all files in root.
	printFolderContent(fss, ".")

	file, err := fss.Open("mail1.template")
	assert.NoError(t, err)

	stat, err := file.Stat()
	assert.NoError(t, err)

	size := stat.Size()
	assert.Greater(t, size, int64(0))

	fdata, err := fs.ReadFile(fss, "mail1.template")
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
