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

// GetFile is for fetching a file from a local file system.
func (sfs *StaticFilesStorage) GetFile(name string) ([]byte, error) {
	filepath, err := model.GetStaticFilePathByFilename(name, sfs.Folder)
	if err != nil {
		return nil, fmt.Errorf("Cannot compose filepath. %s", err)
	}

	if _, err := os.Stat(filepath); err != nil {
		if os.IsNotExist(err) {
			return nil, model.ErrorNotFound
		}
		return nil, fmt.Errorf("Error while checking file '%s' existence. %s", filepath, err)
	}

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("Cannot read %s. %s", filepath, err)
	}
	return data, err
}

// UploadFile is a generic file uploader.
func (sfs *StaticFilesStorage) UploadFile(name string, contents []byte) error {
	filepath, err := model.GetStaticFilePathByFilename(name, sfs.Folder)
	if err != nil {
		return fmt.Errorf("Cannot compose filepath. %s", err)
	}

	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("Cannot open file: %s", err.Error())
	}
	defer file.Close()

	if _, err = io.Copy(file, bytes.NewReader(contents)); err != nil {
		return fmt.Errorf("Cannot save file: %s", err.Error())
	}
	return nil
}

// ParseTemplate parses the html template.
func (sfs *StaticFilesStorage) ParseTemplate(templateName string) (*template.Template, error) {
	filepath, err := model.GetStaticFilePathByFilename(templateName, sfs.Folder)
	if err != nil {
		return nil, fmt.Errorf("Cannot compose filepath. %s", err)
	}
	return template.ParseFiles(filepath)
}

// GetAppleFile is for reading Apple-related static files.
// Unlike generic GetFile, it does not treat model.ErrorNotFound as error.
func (sfs *StaticFilesStorage) GetAppleFile(name string) ([]byte, error) {
	file, err := sfs.GetFile(name)
	if err == model.ErrorNotFound {
		return nil, nil
	}
	return file, err
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
