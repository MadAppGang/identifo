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
	UploadFile(filename Filename, contents io.Reader) error
	StylesHandler() http.Handler
	ScriptsHandler() http.Handler
	ImagesHandler() http.Handler
	FontsHandler() http.Handler
	Close()
}

// StaticPagesNames are the names of html pages.
var StaticPagesNames = StaticPages{
	DisableTFA:            "disable-tfa.html",
	DisableTFASuccess:     "disable-tfa-success.html",
	ForgotPassword:        "forgot-password.html",
	ForgotPasswordSuccess: "forgot-password-success.html",
	Invite:                "invite-email.html",
	Login:                 "login.html",
	Misconfiguration:      "misconfiguration.html",
	Registration:          "registration.html",
	ResetPassword:         "reset-password.html",
	ResetPasswordEmail:    "reset-password-email.html",
	ResetPasswordSuccess:  "reset-password-success.html",
	ResetTFA:              "reset-tfa.html",
	ResetTFASuccess:       "reset-tfa-success.html",
	TokenError:            "token-error.html",
	WebMessage:            "web-message.html",
}

// StaticPages holds together all paths to static pages.
type StaticPages struct {
	DisableTFA            string
	DisableTFASuccess     string
	ForgotPassword        string
	ForgotPasswordSuccess string
	Invite                string
	Login                 string
	Misconfiguration      string
	Registration          string
	ResetPassword         string
	ResetPasswordEmail    string
	ResetPasswordSuccess  string
	ResetTFA              string
	ResetTFASuccess       string
	TFAStatic             string
	TokenError            string
	Verify                string
	WebMessage            string
	Welcome               string
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
	InviteStaticPageName
	LoginStaticPageName
	MisconfigurationStaticPageName
	RegistrationStaticPageName
	ResetPasswordStaticPageName
	ResetPasswordEmailStaticPageName
	ResetPasswordSuccessStaticPageName
	ResetTFAStaticPageName
	ResetTFASuccessStaticPageName
	TFAStaticPageName
	TokenErrorStaticPageName
	VerifyStaticPageName
	WebMessageStaticPageName
	WelcomeStaticPageName
)

// StaticPageName is a name of html template.
type StaticPageName int

// This enum describes names of files related to Apple services.
const (
	AppSiteAssociationFilename Filename = iota + 1
	DeveloperDomainAssociationFilename
)

// AppleFilename is a name of an Apple-related file.
type Filename int
