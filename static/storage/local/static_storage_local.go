package local

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/madappgang/identifo/model"
)

// StaticFilesStorage is a local storage of static files.
type StaticFilesStorage struct {
	staticFilesFolder string
	pagesPath         string
	staticPages       model.StaticPages
	appleFilenames    model.AppleFilenames
}

// NewStaticFilesStorage creates and returns new local static files storage.
func NewStaticFilesStorage(settings model.StaticFilesStorageSettings) (*StaticFilesStorage, error) {
	return &StaticFilesStorage{
		staticFilesFolder: settings.StaticFilesLocation,
		appleFilenames:    settings.AppleFilenames,
	}, nil
}

// ParseTemplate parses the html template.
func (sfs *StaticFilesStorage) ParseTemplate(name model.TemplateName) (*template.Template, error) {
	switch name {
	case model.DisableTFATemplateName:
		return template.ParseFiles(path.Join(sfs.pagesPath, sfs.staticPages.DisableTFA))
	}
	return nil, fmt.Errorf("Unknown template name %v", name)
}

// UploadTemplate is for html template uploads.
func (sfs *StaticFilesStorage) UploadTemplate(templateName model.TemplateName, contents io.Reader) error {
	var filepath string

	switch templateName {
	case model.DisableTFATemplateName:
		filepath = path.Join(sfs.staticFilesFolder, sfs.staticPages.DisableTFA)
	case model.DisableTFASuccessTemplateName:
		filepath = path.Join(sfs.staticFilesFolder, sfs.staticPages.DisableTFASuccess)
	default:
		return fmt.Errorf("Unknown template name %v", templateName)
	}

	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("Cannot open file: %s", err.Error())
	}
	defer file.Close()

	if _, err = io.Copy(file, contents); err != nil {
		return fmt.Errorf("Cannot save file: %s", err.Error())
	}
	return nil
}

// UploadFile is for static file uploads.
func (sfs *StaticFilesStorage) UploadFile(filename model.Filename, contents io.Reader) error {
	var filepath string

	switch filename {
	case model.AppSiteAssociationFilename:
		filepath = path.Join(sfs.staticFilesFolder, sfs.appleFilenames.AppSiteAssociation)
	case model.DeveloperDomainAssociationFilename:
		filepath = path.Join(sfs.staticFilesFolder, sfs.appleFilenames.DeveloperDomainAssociation)
	default:
		return fmt.Errorf("Unknown filename %v", filename)
	}

	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("Cannot open file: %s", err.Error())
	}
	defer file.Close()

	if _, err = io.Copy(file, contents); err != nil {
		return fmt.Errorf("Cannot save file: %s", err.Error())
	}
	return nil
}

// StylesHandler returns server of the css folder.
func (sfs *StaticFilesStorage) StylesHandler() http.Handler {
	stylesHandler := http.FileServer(http.Dir(path.Join(sfs.staticFilesFolder, "/css/")))
	return http.StripPrefix("/css/", stylesHandler)
}

// ScriptsHandler returns server of the js folder.
func (sfs *StaticFilesStorage) ScriptsHandler() http.Handler {
	scriptsHandler := http.FileServer(http.Dir(path.Join(sfs.staticFilesFolder, "/js/")))
	return http.StripPrefix("/js/", scriptsHandler)
}

// ImagesHandler returns server of the img folder.
func (sfs *StaticFilesStorage) ImagesHandler() http.Handler {
	imagesHandler := http.FileServer(http.Dir(path.Join(sfs.staticFilesFolder, "/img/")))
	return http.StripPrefix("/img/", imagesHandler)
}

// FontsHandler returns server of the fonts folder.
func (sfs *StaticFilesStorage) FontsHandler() http.Handler {
	fontsHandler := http.FileServer(http.Dir(path.Join(sfs.staticFilesFolder, "/fonts/")))
	return http.StripPrefix("/fonts/", fontsHandler)
}
