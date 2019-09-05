package model

import (
	"html/template"
	"net/http"
)

// StaticFilesStorage is a wrapper over static files storages.
type StaticFilesStorage interface {
	ParseTemplate(templateName string) (*template.Template, error)
	UploadTemplate(templateName string, contents []byte) error
	ReadAppleFile(filename string) ([]byte, error)
	UploadAppleFile(filename string, contents []byte) error
	AssetHandlers() *AssetHandlers
	AdminPanelHandlers() *AdminPanelHandlers
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
}

// These paths describe directories with static files.
// They are relative to the folder specified in the configuration file.
const (
	AdminPanelBuildPath = "./admin_panel/build"
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
