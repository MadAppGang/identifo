package fs

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

// NewStaticFilesStorage creates and returns new local static files storage.
func NewStaticFilesStorage(settings model.LocalStaticFilesStorageSettings, fallback model.StaticFilesStorage) (model.StaticFilesStorage, error) {
	return &StaticFilesStorage{
		Folder:   settings.FolderPath,
		fallback: fallback,
	}, nil
}

// DefaultFileStorage creates and returns default file storage
func DefaultFileStorage(path string) (model.StaticFilesStorage, error) {
	return &StaticFilesStorage{
		Folder:   path,
		fallback: nil,
	}, nil
}

// StaticFilesStorage is a local storage of static files.
type StaticFilesStorage struct {
	Folder   string
	fallback model.StaticFilesStorage
}

type spaFileSystem struct {
	root http.FileSystem
}

func (fs *spaFileSystem) Open(name string) (http.File, error) {
	f, err := fs.root.Open(name)
	if os.IsNotExist(err) {
		return fs.root.Open("index.html")
	}
	return f, err
}

// GetFile is for fetching a file from a local file system.
func (sfs *StaticFilesStorage) GetFile(name string) ([]byte, error) {
	filepath, err := model.GetStaticFilePathByFilename(name, sfs.Folder)
	if err != nil {
		return nil, fmt.Errorf("Cannot compose filepath. %s", err)
	}

	if _, err := os.Stat(filepath); err != nil {
		if os.IsNotExist(err) {
			if sfs.fallback != nil {
				return sfs.fallback.GetFile(name)
			}
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
		return fmt.Errorf("unable compose the filepath. %s", err)
	}

	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("unable to open the file: %s", err.Error())
	}
	defer file.Close()

	if _, err = io.Copy(file, bytes.NewReader(contents)); err != nil {
		return fmt.Errorf("error saving the file: %s", err.Error())
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

	configHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		jsonData := []byte(`{"apiUrl": "/admin"}`)
		w.Write(jsonData)
	})

	return &model.AdminPanelHandlers{
		SrcHandler:        srcHandler,
		ManagementHandler: managementHandler,
		BuildHandler:      buildHandler,
		ConfigHandler:     configHandler,
	}
}

// WebHandlers returns handlers for the web.
func (sfs *StaticFilesStorage) WebHandlers() *model.WebHandlers {
	appHandler := http.FileServer(&spaFileSystem{http.Dir(path.Join(sfs.Folder, model.WebBuildPath))})

	return &model.WebHandlers{
		AppHandler: appHandler,
	}
}

// Close is to satisfy the interface.
func (sfs *StaticFilesStorage) Close() {}
