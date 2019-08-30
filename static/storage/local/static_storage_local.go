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

var (
	staticPages = model.StaticPagesNames
	appleFiles  = model.AppleFilenames
)

// StaticFilesStorage is a local storage of static files.
type StaticFilesStorage struct {
	staticFilesFolder  string
	pagesPath          string
	emailTemplatesPath string
	appleFilesPath     string
}

// NewStaticFilesStorage creates and returns new local static files storage.
func NewStaticFilesStorage(settings model.StaticFilesStorageSettings) (*StaticFilesStorage, error) {
	return &StaticFilesStorage{
		staticFilesFolder: settings.StaticFilesLocation,
		appleFilesPath:    settings.AppleFilesPath,
	}, nil
}

// ParseTemplate parses the html template.
func (sfs *StaticFilesStorage) ParseTemplate(name model.StaticPageName) (*template.Template, error) {
	switch name {
	case model.DisableTFAStaticPageName:
		return template.ParseFiles(path.Join(sfs.pagesPath, staticPages.DisableTFA))
	case model.DisableTFASuccessStaticPageName:
		return template.ParseFiles(path.Join(sfs.pagesPath, staticPages.DisableTFASuccess))
	case model.ForgotPasswordStaticPageName:
		return template.ParseFiles(path.Join(sfs.pagesPath, staticPages.ForgotPassword))
	case model.ForgotPasswordSuccessStaticPageName:
		return template.ParseFiles(path.Join(sfs.pagesPath, staticPages.ForgotPasswordSuccess))
	case model.InviteStaticPageName:
		return template.ParseFiles(path.Join(sfs.emailTemplatesPath, staticPages.Invite))
	case model.LoginStaticPageName:
		return template.ParseFiles(path.Join(sfs.pagesPath, staticPages.Login))
	case model.MisconfigurationStaticPageName:
		return template.ParseFiles(path.Join(sfs.pagesPath, staticPages.Misconfiguration))
	case model.RegistrationStaticPageName:
		return template.ParseFiles(path.Join(sfs.pagesPath, staticPages.Registration))
	case model.ResetPasswordStaticPageName:
		return template.ParseFiles(path.Join(sfs.pagesPath, staticPages.ResetPassword))
	case model.ResetPasswordEmailStaticPageName:
		return template.ParseFiles(path.Join(sfs.emailTemplatesPath, staticPages.ResetPasswordEmail))
	case model.ResetPasswordSuccessStaticPageName:
		return template.ParseFiles(path.Join(sfs.pagesPath, staticPages.ResetPasswordSuccess))
	case model.ResetTFAStaticPageName:
		return template.ParseFiles(path.Join(sfs.pagesPath, staticPages.ResetTFA))
	case model.ResetTFASuccessStaticPageName:
		return template.ParseFiles(path.Join(sfs.pagesPath, staticPages.ResetTFASuccess))
	case model.TFAStaticPageName:
		return template.ParseFiles(path.Join(sfs.pagesPath, staticPages.TFAStatic))
	case model.TokenErrorStaticPageName:
		return template.ParseFiles(path.Join(sfs.pagesPath, staticPages.TokenError))
	case model.VerifyStaticPageName:
		return template.ParseFiles(path.Join(sfs.pagesPath, staticPages.Verify))
	case model.WebMessageStaticPageName:
		return template.ParseFiles(path.Join(sfs.pagesPath, staticPages.WebMessage))
	case model.WelcomeStaticPageName:
		return template.ParseFiles(path.Join(sfs.pagesPath, staticPages.Welcome))
	}
	return nil, fmt.Errorf("Unsupported template name %v", name)
}

// UploadTemplate is for html template uploads.
func (sfs *StaticFilesStorage) UploadTemplate(templateName model.StaticPageName, contents io.Reader) error {
	var filepath string

	switch templateName {
	case model.DisableTFAStaticPageName:
		filepath = path.Join(sfs.staticFilesFolder, staticPages.DisableTFA)
	case model.DisableTFASuccessStaticPageName:
		filepath = path.Join(sfs.staticFilesFolder, staticPages.DisableTFASuccess)
	case model.ForgotPasswordStaticPageName:
		filepath = path.Join(sfs.pagesPath, staticPages.ForgotPassword)
	case model.ForgotPasswordSuccessStaticPageName:
		filepath = path.Join(sfs.pagesPath, staticPages.ForgotPasswordSuccess)
	case model.InviteStaticPageName:
		filepath = path.Join(sfs.pagesPath, staticPages.Invite)
	case model.LoginStaticPageName:
		filepath = path.Join(sfs.pagesPath, staticPages.Login)
	case model.MisconfigurationStaticPageName:
		filepath = path.Join(sfs.pagesPath, staticPages.Misconfiguration)
	case model.RegistrationStaticPageName:
		filepath = path.Join(sfs.pagesPath, staticPages.Registration)
	case model.ResetPasswordStaticPageName:
		filepath = path.Join(sfs.pagesPath, staticPages.ResetPassword)
	case model.ResetPasswordEmailStaticPageName:
		filepath = path.Join(sfs.pagesPath, staticPages.ResetPasswordEmail)
	case model.ResetPasswordSuccessStaticPageName:
		filepath = path.Join(sfs.pagesPath, staticPages.ResetPasswordSuccess)
	case model.ResetTFAStaticPageName:
		filepath = path.Join(sfs.pagesPath, staticPages.ResetTFA)
	case model.ResetTFASuccessStaticPageName:
		filepath = path.Join(sfs.pagesPath, staticPages.ResetTFASuccess)
	case model.TFAStaticPageName:
		filepath = path.Join(sfs.pagesPath, staticPages.TFAStatic)
	case model.TokenErrorStaticPageName:
		filepath = path.Join(sfs.pagesPath, staticPages.TokenError)
	case model.VerifyStaticPageName:
		filepath = path.Join(sfs.pagesPath, staticPages.Verify)
	case model.WebMessageStaticPageName:
		filepath = path.Join(sfs.pagesPath, staticPages.WebMessage)
	case model.WelcomeStaticPageName:

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
		filepath = path.Join(sfs.staticFilesFolder, appleFiles.AppSiteAssociation)
	case model.DeveloperDomainAssociationFilename:
		filepath = path.Join(sfs.staticFilesFolder, appleFiles.DeveloperDomainAssociation)
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

// Close is to satisfy the interface.
func (sfs *StaticFilesStorage) Close() {}
