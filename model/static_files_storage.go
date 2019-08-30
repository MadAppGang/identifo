package model

import (
	"html/template"
	"io"
	"net/http"
)

// StaticFilesStorage is wrapper over static files storages.
type StaticFilesStorage interface {
	ParseTemplate(templateName TemplateName) (*template.Template, error)
	UploadTemplate(templateName TemplateName, contents io.Reader) error
	UploadFile(filename Filename, contents io.Reader) error
	StylesHandler() http.Handler
	ScriptsHandler() http.Handler
	ImagesHandler() http.Handler
	FontsHandler() http.Handler
}

const (
	DisableTFATemplateName TemplateName = iota + 1
	DisableTFASuccessTemplateName
	ResetTFATemplateName
	LoginTemplateName
	ResetPasswordTemplateName
	ResetPasswordSuccessTemplateName
	WebMessageTemplateName
	RegistrationTemplateName
	ForgotPasswordTemplateName
	ForgotPasswordSuccessTemplateName
	TokenErrorTemplateName
	ResetTFASuccessTemplateName
	MisconfigurationTemplateName
)

// TemplateName is a name of html template.
type TemplateName int

const (
	AppSiteAssociationFilename Filename = iota + 1
	DeveloperDomainAssociationFilename
)

type Filename int
