package model

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strings"
)

// StaticFilesStorage is a wrapper over static files storages.
type StaticFilesStorage interface {
	GetFile(name string) ([]byte, error)
	UploadFile(name string, contents []byte) error
	ParseTemplate(templateName string) (*template.Template, error)
	GetAppleFile(name string) ([]byte, error)
	AssetHandlers() *AssetHandlers
	AdminPanelHandlers() *AdminPanelHandlers
	WebHandlers() *WebHandlers
	Close()
}

// AssetHandlers holds together asset handlers.
type AssetHandlers struct {
	StylesHandler  http.Handler
	ScriptsHandler http.Handler
	ImagesHandler  http.Handler
	FontsHandler   http.Handler
}

// AdminPanelHandlers holds together admin panel handlers.
type AdminPanelHandlers struct {
	SrcHandler        http.Handler
	ManagementHandler http.Handler
	BuildHandler      http.Handler
	ConfigHandler     http.Handler
}

type WebHandlers struct {
	AppHandler http.Handler
}

// These paths describe directories with static files.
// They are relative to the folder specified in the configuration file.
const (
	AdminPanelBuildPath = "./admin_panel"
	WebBuildPath        = "./web"
	PagesPath           = "./html"
	EmailTemplatesPath  = "./email_templates"
	AppleFilesPath      = "./apple"
)

// StaticPagesNames are the names of html pages.
var StaticPagesNames = StaticPages{
	DisableTFA:            "disable-tfa.html",
	DisableTFASuccess:     "disable-tfa-success.html",
	ForgotPassword:        "forgot-password.html",
	ForgotPasswordSuccess: "forgot-password-success.html",
	InviteEmail:           "invite-email.html",
	Login:                 "login.html",
	Misconfiguration:      "misconfiguration.html",
	Registration:          "registration.html",
	ResetPassword:         "reset-password.html",
	ResetPasswordEmail:    "reset-password-email.html",
	ResetPasswordSuccess:  "reset-password-success.html",
	ResetTFA:              "reset-tfa.html",
	ResetTFASuccess:       "reset-tfa-success.html",
	TFAEmail:              "tfa-email.html",
	TokenError:            "token-error.html",
	VerifyEmail:           "verify-email.html",
	WebMessage:            "web-message.html",
	WelcomeEmail:          "welcome-email.html",
}

// StaticPages holds together all paths to static pages.
type StaticPages struct {
	DisableTFA            string
	DisableTFASuccess     string
	ForgotPassword        string
	ForgotPasswordSuccess string
	InviteEmail           string
	Login                 string
	Misconfiguration      string
	Registration          string
	ResetPassword         string
	ResetPasswordEmail    string
	ResetPasswordSuccess  string
	ResetTFA              string
	ResetTFASuccess       string
	TFAEmail              string
	TokenError            string
	VerifyEmail           string
	WebMessage            string
	WelcomeEmail          string
}

// AppleFiles holds together static files needed for supporting Apple services.
type AppleFiles struct {
	DeveloperDomainAssociation string `yaml:"developerDomainAssociation,omitempty" json:"developer_domain_association,omitempty"`
	AppSiteAssociation         string `yaml:"appSiteAssociation,omitempty" json:"app_site_association,omitempty"`
}

// AppleFilenames are names of the files related to Apple services.
var AppleFilenames = AppleFiles{
	DeveloperDomainAssociation: "apple-developer-domain-association.txt",
	AppSiteAssociation:         "apple-app-site-association",
}

// GetStaticFilePathByFilename returns filepath for given static file name.
func GetStaticFilePathByFilename(filename, staticFolder string) (filepath string, err error) {
	if strings.Contains(filename, "apple") {
		return path.Join(staticFolder, AppleFilesPath, filename), nil
	}

	switch path.Ext(filename) {
	case ".html":
		if strings.Contains(filename, "email") {
			filepath = path.Join(staticFolder, EmailTemplatesPath, filename)
		} else {
			filepath = path.Join(staticFolder, PagesPath, filename)
		}
	case ".css":
		filepath = path.Join(staticFolder, "css", filename)
	case ".js":
		filepath = path.Join(staticFolder, "js/dist", filename)
	case ".png", ".jpg", ".jpeg":
		filepath = path.Join(staticFolder, "img", filename)
	case ".woff":
		filepath = path.Join(staticFolder, "fonts", filename)
	default:
		err = fmt.Errorf("Unknown extension '%s'", path.Ext(filename))
	}
	return
}
