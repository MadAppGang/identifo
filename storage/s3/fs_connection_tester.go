package s3

import (
	"fmt"
	"io/fs"

	"github.com/madappgang/identifo/model"
)

func NewFileStorageConnectionTester(settings model.FileStorageS3, expectedFiles []string) model.ConnectionTester {
	return &S3FSConnectionTester{
		settings:      settings,
		expectedFiles: expectedFiles,
	}
}

type S3FSConnectionTester struct {
	settings      model.FileStorageS3
	expectedFiles []string
}

func (ct *S3FSConnectionTester) Connect() error {
	s3fs, err := NewFS(ct.settings)
	if err != nil {
		return err
	}

	if len(ct.expectedFiles) == 0 {
		// if we not expecting any files, let's just try to read the contnt of current folder
		_, err := fs.ReadDir(s3fs, "./")
		return err
	}

	for _, f := range ct.expectedFiles {
		fs, err := s3fs.Open(f)
		if err != nil {
			return fmt.Errorf("error to get expected file: %s, error: %v", f, err)
		}
		fs.Close()
	}

	return nil
}
