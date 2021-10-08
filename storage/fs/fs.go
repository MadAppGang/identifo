package fs

import (
	"io/fs"

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
