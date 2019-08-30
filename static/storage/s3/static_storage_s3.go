package s3

import (
	"html/template"
	"io"
	"net/http"

	"github.com/madappgang/identifo/model"
)

// StaticFilesStorage is a storage of static files in S3.
type StaticFilesStorage struct {
	bucket             string
	pagesPath          string
	emailTemplatesPath string
}

// NewStaticFilesStorage creates and returns new static files storage in S3.
func NewStaticFilesStorage(settings model.StaticFilesStorageSettings) (*StaticFilesStorage, error) {
	return &StaticFilesStorage{
		bucket:    settings.StaticFilesLocation,
		pagesPath: settings.PagesPath,
	}, nil
}

// ParseTemplate parses the html template.
func (sfs *StaticFilesStorage) ParseTemplate(templateName model.StaticPageName) (*template.Template, error) {
	// TODO: implement
	return nil, nil
}

// UploadTemplate is for html template uploads.
func (sfs *StaticFilesStorage) UploadTemplate(templateName model.StaticPageName, contents io.Reader) error {
	// TODO: implement
	return nil
}

// UploadFile is for static file uploads.
func (sfs *StaticFilesStorage) UploadFile(filename model.Filename, contents io.Reader) error {
	// TODO: implement
	return nil
}

// StylesHandler returns server of the images folder.
func (sfs *StaticFilesStorage) StylesHandler() http.Handler {
	// TODO: implement
	return nil
}

// ScriptsHandler returns server of the images folder.
func (sfs *StaticFilesStorage) ScriptsHandler() http.Handler {
	// TODO: implement
	return nil
}

// ImagesHandler returns server of the images folder.
func (sfs *StaticFilesStorage) ImagesHandler() http.Handler {
	// TODO: implement
	return nil
}

// FontsHandler returns server of the fonts folder.
func (sfs *StaticFilesStorage) FontsHandler() http.Handler {
	// TODO: implement
	return nil
}

// Close is to satisfy the interface.
func (sfs *StaticFilesStorage) Close() {}
