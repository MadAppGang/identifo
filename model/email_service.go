package model

import (
	"html/template"
	"path"
)

// EmailService manage sending email
type EmailService interface {
	SendMessage(subject, body, recipient string) error
	SendHTML(subject, html, recipient string) error

	SendTemplateEmail(subject, recipient string, template *template.Template, data interface{}) error
	SendResetEmail(subject, recipient string, data interface{}) error
	SendWelcomeEmail(subject, recipient string, data interface{}) error
	SendVerifyEmail(subject, recipient string, data interface{}) error

	Templater() *EmailTemplater
}

//DefaultEmailTemplater creates and returns new email templater with default settings
func DefaultEmailTemplater() (*EmailTemplater, error) {
	return NewEmailTemplater(DefaultEmailTemplates, "./email_templates")
}

//NewEmailTemplater creates new templater
func NewEmailTemplater(templates EmailTemplates, templatePath string) (*EmailTemplater, error) {
	et := EmailTemplater{}
	var err error
	f := path.Join(templatePath, templates.Welcome)
	if et.WelcomeTemplate, err = template.New(path.Base(f)).ParseFiles(f); err != nil {
		return nil, err
	}
	f = path.Join(templatePath, templates.ResetPassword)
	if et.ResetPasswordTemplate, err = template.New(path.Base(f)).ParseFiles(f); err != nil {
		return nil, err
	}
	f = path.Join(templatePath, templates.VerifyEmail)
	if et.VerifyEmailTemplate, err = template.New(path.Base(f)).ParseFiles(f); err != nil {
		return nil, err
	}
	return &et, nil
}

//EmailTemplater creates and stores email templates
type EmailTemplater struct {
	WelcomeTemplate       *template.Template
	ResetPasswordTemplate *template.Template
	VerifyEmailTemplate   *template.Template
}

//EmailTemplates store email templates
type EmailTemplates struct {
	Welcome       string
	ResetPassword string
	VerifyEmail   string
}

//DefaultEmailTemplates stores default email template names
var DefaultEmailTemplates = EmailTemplates{
	Welcome:       "welcome.html",
	ResetPassword: "reset_password.html",
	VerifyEmail:   "verify_email.html",
}
