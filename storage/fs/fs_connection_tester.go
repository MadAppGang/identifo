package fs

import (
	"fmt"
	"io/fs"

	"github.com/madappgang/identifo/v2/model"
)

func NewFSConnectionTester(settings model.FileStorageLocal, expectedFiles []string) model.ConnectionTester {
	return &FSConnectionTester{
		settings:      settings,
		expectedFiles: expectedFiles,
	}
}

type FSConnectionTester struct {
	settings      model.FileStorageLocal
	expectedFiles []string
}

func (ct *FSConnectionTester) Connect() error {
	fss := NewFS(ct.settings)
	if len(ct.expectedFiles) == 0 {
		// if we not expecting any files, let's just try to read the contnt of current folder
		_, err := fs.ReadDir(fss, "./")
		return err
	}

	for _, f := range ct.expectedFiles {
		fs, err := fss.Open(f)
		if err != nil {
			return fmt.Errorf("error to get expected file: %s, error: %v", f, err)
		}
		fs.Close()
	}
	return nil
}
