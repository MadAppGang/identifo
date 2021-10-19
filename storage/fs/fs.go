package fs

import (
	"io/fs"
	"path"

	"github.com/madappgang/identifo/model"
	"github.com/spf13/afero"
)

func NewFS(settings model.FileStorageLocal) fs.FS {
	// using afero file system abstraction to provice the FS chaining:
	// - io.FS wrapper
	// - BasePathFS wrapper
	// - regular os.fs wrapper

	return afero.NewIOFS(
		afero.NewBasePathFs(
			afero.NewOsFs(),
			settings.FolderPath,
		),
	)
}

// NewFSWithFiles creates the fs which already has predefined files on top of the base fs
func NewFSWithFiles(settings model.FileStorageLocal, files map[string][]byte) fs.FS {
	base := afero.NewBasePathFs(afero.NewOsFs(), settings.FolderPath)
	layer := afero.NewMemMapFs()
	for filename, data := range files {
		afero.WriteFile(layer, path.Join(settings.FolderPath, filename), data, 0644)
	}
	combined := afero.NewCopyOnWriteFs(base, layer)

	return afero.NewIOFS(
		afero.NewBasePathFs(
			combined,
			settings.FolderPath,
		),
	)
}
