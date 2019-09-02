package s3

import (
	"html/template"
	"io"

	"github.com/madappgang/identifo/model"
)

// StaticFilesStorage is a storage of static files in S3.
type StaticFilesStorage struct {
	bucket             string
	pagesPath          string
	emailTemplatesPath string
	appleFilesPath     string
}

// NewStaticFilesStorage creates and returns new static files storage in S3.
func NewStaticFilesStorage(settings model.StaticFilesStorageSettings) (*StaticFilesStorage, error) {
	return &StaticFilesStorage{
		bucket:             settings.StaticFilesLocation,
		pagesPath:          settings.PagesPath,
		emailTemplatesPath: settings.EmailTemplatesPath,
		appleFilesPath:     settings.AppleFilesPath,
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

// ReadAppleFile is for reading Apple-related static files.
func (sfs *StaticFilesStorage) ReadAppleFile(filename model.AppleFilename) ([]byte, error) {
	// TODO: implement
	return nil, nil
}

// UploadAppleFile is for Apple-related file uploads.
func (sfs *StaticFilesStorage) UploadAppleFile(filename model.AppleFilename, contents io.Reader) error {
	// TODO: implement
	return nil
}

// AssetHandlers returns handlers for assets.
func (sfs *StaticFilesStorage) AssetHandlers() *model.AssetHandlers {
	// TODO: implement
	return nil
}

// AdminPanelHandlers returns handlers for the admin panel.
func (sfs *StaticFilesStorage) AdminPanelHandlers() *model.AdminPanelHandlers {
	// TODO: implement
	return nil
}

// Close is to satisfy the interface.
func (sfs *StaticFilesStorage) Close() {}
