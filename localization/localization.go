package localization

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gopkg.in/yaml.v2"
)

//go:embed translations/*.yaml
var translationFS embed.FS

// SupportedLangs the list of currently supported languages
var SupportedLangs = []language.Tag{}

// LoadDefaultCatalog loads data from default catalog
func LoadDefaultCatalog() error {
	// check if the data is already being loaded
	if len(SupportedLangs) != 0 {
		return nil
	}

	err := fs.WalkDir(translationFS, "translations", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// skip dirs
		if d.IsDir() {
			return nil
		}

		fd, err := fs.ReadFile(translationFS, p)
		if err != nil {
			return nil
		}

		tagStr := strings.TrimSuffix(d.Name(), filepath.Ext(d.Name()))
		tag, err := language.Parse(tagStr)
		if err != nil {
			fmt.Printf("unable to load translation file %s, as could not construct language tag from file name with error: %s\n", d.Name(), err.Error())
			// not returning error, just moving forward to other files to try to load other languages
			return nil
		}

		data := map[string]string{}
		err = yaml.Unmarshal(fd, &data)
		if err != nil {
			return err
		}

		for k, v := range data {
			message.SetString(tag, k, v)
		}

		return nil
	})
	return err
}
