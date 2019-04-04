package model

import (
	"html/template"
	"path"
)

// EmailService manages sending emails.
type EmailService interface {
	SendMessage(subject, body, recipient string) error
	SendHTML(subject, html, recipient string) error

	SendTemplateEmail(subject, recipient string, template *template.Template, data interface{}) error
	SendResetEmail(subject, recipient string, data interface{}) error
	SendWelcomeEmail(subject, recipient string, data interface{}) error
	SendVerifyEmail(subject, recipient string, data interface{}) error

	Templater() *EmailTemplater
}

// NewEmailTemplater creates new email templater.
func NewEmailTemplater(templateNames EmailTemplateNames, templatePath string) (*EmailTemplater, error) {
	et := EmailTemplater{}
	var err error

	f := path.Join(templatePath, templateNames.Welcome)
	if et.WelcomeTemplate, err = template.New(path.Base(f)).ParseFiles(f); err != nil {
		return nil, err
	}

	f = path.Join(templatePath, templateNames.ResetPassword)
	if et.ResetPasswordTemplate, err = template.New(path.Base(f)).ParseFiles(f); err != nil {
		return nil, err
	}

	f = path.Join(templatePath, templateNames.VerifyEmail)
	if et.VerifyEmailTemplate, err = template.New(path.Base(f)).ParseFiles(f); err != nil {
		return nil, err
	}
	return &et, nil
}

// EmailTemplater stores pointers to email templates.
type EmailTemplater struct {
	WelcomeTemplate       *template.Template
	ResetPasswordTemplate *template.Template
	VerifyEmailTemplate   *template.Template
}

// EmailTemplateNames stores email template names.
type EmailTemplateNames struct {
	Welcome       string `yaml:"welcome,omitempty" json:"welcome,omitempty"`
	ResetPassword string `yaml:"resetPassword,omitempty" json:"reset_password,omitempty"`
	VerifyEmail   string `yaml:"verifyEmail,omitempty" json:"verify_email,omitempty"`
}
