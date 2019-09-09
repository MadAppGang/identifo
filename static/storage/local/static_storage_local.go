package local

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/madappgang/identifo/model"
)

// StaticFilesStorage is a local storage of static files.
type StaticFilesStorage struct {
	Folder string
}

// NewStaticFilesStorage creates and returns new local static files storage.
func NewStaticFilesStorage(settings model.StaticFilesStorageSettings) (*StaticFilesStorage, error) {
	return &StaticFilesStorage{
		Folder: settings.Folder,
	}, nil
}

// ParseTemplate parses the html template.
func (sfs *StaticFilesStorage) ParseTemplate(templateName string) (*template.Template, error) {
	if strings.Contains(strings.ToLower(templateName), "email") {
		templateName = path.Join(sfs.Folder, model.EmailTemplatesPath, templateName)
	} else {
		templateName = path.Join(sfs.Folder, model.PagesPath, templateName)
	}
	return template.ParseFiles(templateName)
}

// UploadTemplate is for html template uploads.
func (sfs *StaticFilesStorage) UploadTemplate(templateName string, contents []byte) error {
	if strings.Contains(strings.ToLower(templateName), "email") {
		templateName = path.Join(sfs.Folder, model.EmailTemplatesPath, templateName)
	} else {
		templateName = path.Join(sfs.Folder, model.PagesPath, templateName)
	}

	file, err := os.OpenFile(templateName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("Cannot open file: %s", err.Error())
	}
	defer file.Close()

	if _, err = io.Copy(file, bytes.NewReader(contents)); err != nil {
		return fmt.Errorf("Cannot save file: %s", err.Error())
	}
	return nil
}

// ReadAppleFile is for reading Apple-related static files.
func (sfs *StaticFilesStorage) ReadAppleFile(filename string) ([]byte, error) {
	filename = path.Join(sfs.Folder, model.AppleFilesPath, filename)
	// Check if file exists. If not - return nil error and nil slice.
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("Error while checking filename '%s' existence. %s", filename, err)
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Cannot read %s. %s", filename, err)
	}
	return data, nil
}

// UploadAppleFile is for Apple-related file uploads.
func (sfs *StaticFilesStorage) UploadAppleFile(filename string, contents []byte) error {
	filename = path.Join(sfs.Folder, model.AppleFilesPath, filename)

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("Cannot open file: %s", err.Error())
	}
	defer file.Close()

	if _, err = io.Copy(file, bytes.NewReader(contents)); err != nil {
		return fmt.Errorf("Cannot save file: %s", err.Error())
	}
	return nil
}

// AssetHandlers returns handlers for assets.
func (sfs *StaticFilesStorage) AssetHandlers() *model.AssetHandlers {
	stylesHandler := http.FileServer(http.Dir(path.Join(sfs.Folder, "/css/")))
	scriptsHandler := http.FileServer(http.Dir(path.Join(sfs.Folder, "/js/")))
	imagesHandler := http.FileServer(http.Dir(path.Join(sfs.Folder, "/img/")))
	fontsHandler := http.FileServer(http.Dir(path.Join(sfs.Folder, "/fonts/")))

	return &model.AssetHandlers{
		StylesHandler:  http.StripPrefix("/css/", stylesHandler),
		ScriptsHandler: http.StripPrefix("/js/", scriptsHandler),
		ImagesHandler:  http.StripPrefix("/img/", imagesHandler),
		FontsHandler:   http.StripPrefix("/fonts/", fontsHandler),
	}
}

// AdminPanelHandlers returns handlers for the admin panel.
func (sfs *StaticFilesStorage) AdminPanelHandlers() *model.AdminPanelHandlers {
	srcHandler := http.StripPrefix("/src/", http.FileServer(http.Dir(path.Join(sfs.Folder, model.AdminPanelBuildPath, "/src"))))

	managementHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join(sfs.Folder, model.AdminPanelBuildPath, "/index.html"))
	})

	buildHandler := http.FileServer(http.Dir(path.Join(sfs.Folder, model.AdminPanelBuildPath)))

	return &model.AdminPanelHandlers{
		SrcHandler:        srcHandler,
		ManagementHandler: managementHandler,
		BuildHandler:      buildHandler,
	}
}

// Close is to satisfy the interface.
func (sfs *StaticFilesStorage) Close() {}
