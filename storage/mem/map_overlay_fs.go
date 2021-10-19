package mem

import (
	"io/fs"
	"os"
	"path"

	"github.com/spf13/afero"
)

func NewMapOverlayFS(base fs.FS, files map[string][]byte) fs.FS {
	layer := afero.NewMemMapFs()
	for filename, data := range files {
		afero.WriteFile(layer, path.Clean(filename), data, 0644)
	}

	return &mapFS{
		base:  base,
		memFS: afero.NewIOFS(layer),
	}
}

type mapFS struct {
	base  fs.FS
	memFS fs.FS
}

func (f *mapFS) Open(name string) (fs.File, error) {
	mf, err := f.memFS.Open(name)
	if err != nil {
		if os.IsNotExist(err) || fs.ErrInvalid == err {
			return f.base.Open(name)
		}
		return nil, err
	}
	return mf, nil
}
