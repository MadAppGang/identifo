package s3

import (
	"html/template"
	"io"
	"net/http"

	"github.com/madappgang/identifo/model"
)

// StaticFilesStorage is a storage of static files in S3.
type StaticFilesStorage struct {
	bucket         string
	pagesPath      string
	staticPages    model.StaticPages
	appleFilenames model.AppleFilenames
}

// NewStaticFilesStorage creates and returns new static files storage in S3.
func NewStaticFilesStorage(settings model.StaticFilesStorageSettings) (*StaticFilesStorage, error) {
	return &StaticFilesStorage{
		bucket:         settings.StaticFilesLocation,
		appleFilenames: settings.AppleFilenames,
	}, nil
}

func (sfs *StaticFilesStorage) ParseTemplate(templateName model.TemplateName) (*template.Template, error) {
	// TODO: implement
	return nil, nil
}

func (sfs *StaticFilesStorage) UploadTemplate(templateName model.TemplateName, contents io.Reader) error {
	// TODO: implement
	return nil
}

func (sfs *StaticFilesStorage) UploadFile(filename model.Filename, contents io.Reader) error {
	// TODO: implement
	return nil
}

func (sfs *StaticFilesStorage) StylesHandler() http.Handler {
	// TODO: implement
	return nil
}

func (sfs *StaticFilesStorage) ScriptsHandler() http.Handler {
	// TODO: implement
	return nil
}

func (sfs *StaticFilesStorage) ImagesHandler() http.Handler {
	// TODO: implement
	return nil
}

func (sfs *StaticFilesStorage) FontsHandler() http.Handler {
	// TODO: implement
	return nil
}
