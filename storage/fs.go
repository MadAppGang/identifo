package storage

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/madappgang/identifo/v2/model"
	fss "github.com/madappgang/identifo/v2/storage/fs"
	"github.com/madappgang/identifo/v2/storage/s3"
)

func NewFS(settings model.FileStorageSettings) (fs.FS, error) {
	switch settings.Type {
	case model.FileStorageTypeLocal:
		ss := fss.NewFS(settings.Local)
		return ss, nil
		// return &RootReplacedFS{Root: settings.Local.FolderPath, FS: ss}, nil
	case model.FileStorageTypeS3:
		ss, err := s3.NewFS(settings.S3)
		if err != nil {
			return nil, err
		}
		return &RootReplacedFS{Root: settings.S3.Key, FS: ss}, nil
	default:
		return nil, fmt.Errorf("file storage type is not supported %s ", settings.Type)
	}
}

// RootReplacedFS add root prefix on top of every url of underlying fs.FS
type RootReplacedFS struct {
	Root string
	FS   fs.FS
}

// we add root path every time we want to open the file
func (f *RootReplacedFS) Open(name string) (fs.File, error) {
	fn := filepath.Join(f.Root, filepath.Clean("/"+name))
	return f.FS.Open(fn)
}
