package model

import (
	"html/template"
)

const (
	DisableTFATemplateName TemplateName = "disable_tfa"
	ResetTFATemplateName   TemplateName = "reset_tfa"
	LoginTemplateName      TemplateName = "login"
)

// TemplateName is a name of html template.
type TemplateName string

// StaticFilesStorage is wrapper over static files storages.
type StaticFilesStorage interface {
	ParseTemplate(name TemplateName) (*template.Template, error)
	UploadFile(filename string) error
}
