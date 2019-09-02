package model

import (
	"html/template"
	"io"
	"net/http"
)

// StaticFilesStorage is a wrapper over static files storages.
type StaticFilesStorage interface {
	ParseTemplate(templateName StaticPageName) (*template.Template, error)
	UploadTemplate(templateName StaticPageName, contents io.Reader) error
	ReadAppleFile(filename AppleFilename) ([]byte, error)
	UploadAppleFile(filename AppleFilename, contents io.Reader) error
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

// This enum describes static pages names.
const (
	DisableTFAStaticPageName StaticPageName = iota + 1
	DisableTFASuccessStaticPageName
	ForgotPasswordStaticPageName
	ForgotPasswordSuccessStaticPageName
	InviteEmailStaticPageName
	LoginStaticPageName
	MisconfigurationStaticPageName
	RegistrationStaticPageName
	ResetPasswordStaticPageName
	ResetPasswordEmailStaticPageName
	ResetPasswordSuccessStaticPageName
	ResetTFAStaticPageName
	ResetTFASuccessStaticPageName
	TFAEmailStaticPageName
	TokenErrorStaticPageName
	VerifyEmailStaticPageName
	WebMessageStaticPageName
	WelcomeEmailStaticPageName
)

// StaticPageName is a name of html template.
type StaticPageName int

// This enum describes names of files related to Apple services.
const (
	AppSiteAssociationFilename AppleFilename = iota + 1
	DeveloperDomainAssociationFilename
)

// AppleFilename is a name of an Apple-related file.
type AppleFilename int
