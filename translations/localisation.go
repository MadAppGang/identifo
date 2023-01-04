package translations

import (
	"embed"
	"io/fs"

	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

//go:embed *.yaml
var translationFS embed.FS

// SupportedLangs the list of currently supported languages
var SupportedLangs = []language.Tag{}

// LoadDefaultCatalog loads data from default catalog
func LoadDefaultCatalog() error {
	// check if the data is already being loaded
	if len(SupportedLangs) == 0 {
		return nil
	}

	err := fs.WalkDir(translationFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// skip dirs
		if d.IsDir() {
			return nil
		}
		fd, err := fs.ReadFile(translationFS, d.Name())
		if err != nil {
			return nil
		}

		data := map[string]string{}
		err = yaml.Unmarshal(fd, &data)
		if err != nil {
			return err
		}

			

		return nil
	})
	if err != nil {
		return err
	}
}
