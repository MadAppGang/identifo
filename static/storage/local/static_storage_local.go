package local

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
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
	staticFilesFolder   string
	pagesPath           string
	emailTemplatesPath  string
	adminPanelBuildPath string
	appleFilesPath      string
}

// NewStaticFilesStorage creates and returns new local static files storage.
func NewStaticFilesStorage(settings model.StaticFilesStorageSettings) (*StaticFilesStorage, error) {
	return &StaticFilesStorage{
		staticFilesFolder:   settings.StaticFilesLocation,
		pagesPath:           settings.PagesPath,
		emailTemplatesPath:  settings.EmailTemplatesPath,
		adminPanelBuildPath: settings.AdminPanelBuildPath,
		appleFilesPath:      settings.AppleFilesPath,
	}, nil
}

// ParseTemplate parses the html template.
func (sfs *StaticFilesStorage) ParseTemplate(name model.StaticPageName) (*template.Template, error) {
	pagesPath := path.Join(sfs.staticFilesFolder, sfs.pagesPath)
	emailsPath := path.Join(sfs.staticFilesFolder, sfs.emailTemplatesPath)

	switch name {
	case model.DisableTFAStaticPageName:
		return template.ParseFiles(path.Join(pagesPath, staticPages.DisableTFA))
	case model.DisableTFASuccessStaticPageName:
		return template.ParseFiles(path.Join(pagesPath, staticPages.DisableTFASuccess))
	case model.ForgotPasswordStaticPageName:
		return template.ParseFiles(path.Join(pagesPath, staticPages.ForgotPassword))
	case model.ForgotPasswordSuccessStaticPageName:
		return template.ParseFiles(path.Join(pagesPath, staticPages.ForgotPasswordSuccess))
	case model.InviteEmailStaticPageName:
		return template.ParseFiles(path.Join(emailsPath, staticPages.InviteEmail))
	case model.LoginStaticPageName:
		return template.ParseFiles(path.Join(pagesPath, staticPages.Login))
	case model.MisconfigurationStaticPageName:
		return template.ParseFiles(path.Join(pagesPath, staticPages.Misconfiguration))
	case model.RegistrationStaticPageName:
		return template.ParseFiles(path.Join(pagesPath, staticPages.Registration))
	case model.ResetPasswordStaticPageName:
		return template.ParseFiles(path.Join(pagesPath, staticPages.ResetPassword))
	case model.ResetPasswordEmailStaticPageName:
		return template.ParseFiles(path.Join(emailsPath, staticPages.ResetPasswordEmail))
	case model.ResetPasswordSuccessStaticPageName:
		return template.ParseFiles(path.Join(pagesPath, staticPages.ResetPasswordSuccess))
	case model.ResetTFAStaticPageName:
		return template.ParseFiles(path.Join(pagesPath, staticPages.ResetTFA))
	case model.ResetTFASuccessStaticPageName:
		return template.ParseFiles(path.Join(pagesPath, staticPages.ResetTFASuccess))
	case model.TFAEmailStaticPageName:
		return template.ParseFiles(path.Join(emailsPath, staticPages.TFAEmail))
	case model.TokenErrorStaticPageName:
		return template.ParseFiles(path.Join(pagesPath, staticPages.TokenError))
	case model.VerifyEmailStaticPageName:
		return template.ParseFiles(path.Join(emailsPath, staticPages.VerifyEmail))
	case model.WebMessageStaticPageName:
		return template.ParseFiles(path.Join(pagesPath, staticPages.WebMessage))
	case model.WelcomeEmailStaticPageName:
		return template.ParseFiles(path.Join(emailsPath, staticPages.WelcomeEmail))
	}
	return nil, fmt.Errorf("Unsupported template name %v", name)
}

// UploadTemplate is for html template uploads.
func (sfs *StaticFilesStorage) UploadTemplate(templateName model.StaticPageName, contents io.Reader) error {
	pagesPath := path.Join(sfs.staticFilesFolder, sfs.pagesPath)
	emailsPath := path.Join(sfs.staticFilesFolder, sfs.emailTemplatesPath)

	var filepath string

	switch templateName {
	case model.DisableTFAStaticPageName:
		filepath = path.Join(pagesPath, staticPages.DisableTFA)
	case model.DisableTFASuccessStaticPageName:
		filepath = path.Join(pagesPath, staticPages.DisableTFASuccess)
	case model.ForgotPasswordStaticPageName:
		filepath = path.Join(pagesPath, staticPages.ForgotPassword)
	case model.ForgotPasswordSuccessStaticPageName:
		filepath = path.Join(pagesPath, staticPages.ForgotPasswordSuccess)
	case model.InviteEmailStaticPageName:
		filepath = path.Join(emailsPath, staticPages.InviteEmail)
	case model.LoginStaticPageName:
		filepath = path.Join(pagesPath, staticPages.Login)
	case model.MisconfigurationStaticPageName:
		filepath = path.Join(pagesPath, staticPages.Misconfiguration)
	case model.RegistrationStaticPageName:
		filepath = path.Join(pagesPath, staticPages.Registration)
	case model.ResetPasswordStaticPageName:
		filepath = path.Join(pagesPath, staticPages.ResetPassword)
	case model.ResetPasswordEmailStaticPageName:
		filepath = path.Join(emailsPath, staticPages.ResetPasswordEmail)
	case model.ResetPasswordSuccessStaticPageName:
		filepath = path.Join(pagesPath, staticPages.ResetPasswordSuccess)
	case model.ResetTFAStaticPageName:
		filepath = path.Join(pagesPath, staticPages.ResetTFA)
	case model.ResetTFASuccessStaticPageName:
		filepath = path.Join(pagesPath, staticPages.ResetTFASuccess)
	case model.TFAEmailStaticPageName:
		filepath = path.Join(emailsPath, staticPages.TFAEmail)
	case model.TokenErrorStaticPageName:
		filepath = path.Join(pagesPath, staticPages.TokenError)
	case model.VerifyEmailStaticPageName:
		filepath = path.Join(emailsPath, staticPages.VerifyEmail)
	case model.WebMessageStaticPageName:
		filepath = path.Join(pagesPath, staticPages.WebMessage)
	case model.WelcomeEmailStaticPageName:
		filepath = path.Join(emailsPath, staticPages.WelcomeEmail)
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

// ReadAppleFile is for reading Apple-related static files.
func (sfs *StaticFilesStorage) ReadAppleFile(filename model.AppleFilename) ([]byte, error) {
	appleFolder := path.Join(sfs.staticFilesFolder, sfs.appleFilesPath)
	var filepath string

	switch filename {
	case model.AppSiteAssociationFilename:
		filepath = path.Join(appleFolder, appleFiles.AppSiteAssociation)
	case model.DeveloperDomainAssociationFilename:
		filepath = path.Join(appleFolder, appleFiles.DeveloperDomainAssociation)
	default:
		return nil, fmt.Errorf("Unknown filename %v", filename)
	}

	// Check if file exists. If not - return nil error and nil slice.
	if _, err := os.Stat(filepath); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("Error while checking filepath '%s' existence. %s", filepath, err)
	}

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("Cannot read %s. %s", filepath, err)
	}
	return data, nil
}

// UploadAppleFile is for Apple-related file uploads.
func (sfs *StaticFilesStorage) UploadAppleFile(filename model.AppleFilename, contents io.Reader) error {
	appleFolder := path.Join(sfs.staticFilesFolder, sfs.appleFilesPath)
	var filepath string

	switch filename {
	case model.AppSiteAssociationFilename:
		filepath = path.Join(appleFolder, appleFiles.AppSiteAssociation)
	case model.DeveloperDomainAssociationFilename:
		filepath = path.Join(appleFolder, appleFiles.DeveloperDomainAssociation)
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

// AssetHandlers returns handlers for assets.
func (sfs *StaticFilesStorage) AssetHandlers() *model.AssetHandlers {
	stylesHandler := http.FileServer(http.Dir(path.Join(sfs.staticFilesFolder, "/css/")))
	scriptsHandler := http.FileServer(http.Dir(path.Join(sfs.staticFilesFolder, "/js/")))
	imagesHandler := http.FileServer(http.Dir(path.Join(sfs.staticFilesFolder, "/img/")))
	fontsHandler := http.FileServer(http.Dir(path.Join(sfs.staticFilesFolder, "/fonts/")))

	return &model.AssetHandlers{
		StylesHandler:  http.StripPrefix("/css/", stylesHandler),
		ScriptsHandler: http.StripPrefix("/js/", scriptsHandler),
		ImagesHandler:  http.StripPrefix("/img/", imagesHandler),
		FontsHandler:   http.StripPrefix("/fonts/", fontsHandler),
	}
}

// AdminPanelHandlers returns handlers for the admin panel.
func (sfs *StaticFilesStorage) AdminPanelHandlers() *model.AdminPanelHandlers {
	adminPanelFolder := path.Join(sfs.staticFilesFolder, sfs.adminPanelBuildPath)

	srcHandler := http.StripPrefix("/src/", http.FileServer(http.Dir(path.Join(adminPanelFolder, "/src"))))
	managementHandleFunc := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join(adminPanelFolder, "/index.html"))
	}
	buildHandler := http.FileServer(http.Dir(adminPanelFolder))

	return &model.AdminPanelHandlers{
		SrcHandler:        srcHandler,
		ManagementHandler: http.HandlerFunc(managementHandleFunc),
		BuildHandler:      buildHandler,
	}
}

// Close is to satisfy the interface.
func (sfs *StaticFilesStorage) Close() {}
