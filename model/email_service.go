package model

import (
	"html/template"
)

// EmailService manages sending emails.
type EmailService interface {
	SendMessage(subject, body, recipient string) error
	SendHTML(subject, html, recipient string) error

	SendTemplateEmail(subject, recipient string, template *template.Template, data interface{}) error
	SendResetEmail(subject, recipient string, data ResetEmailData) error
	SendInviteEmail(subject, recipient string, data InviteEmailData) error
	SendWelcomeEmail(subject, recipient string, data WelcomeEmailData) error
	SendVerifyEmail(subject, recipient string, data VerifyEmailData) error
	SendTFAEmail(subject, recipient string, data SendTFAEmailData) error

	Templater() *EmailTemplater
}

// ResetEmailData represents data to be send to the user for reset email
type ResetEmailData struct {
	User  User
	Token string
	URL   string
	Host  string
	Data  interface{}
}

type InviteEmailData struct {
	Requester User
	Token     string
	URL       string
	Host      string
	Query     string
	App       string
	Scopes    string
	Callback  string
	Data      interface{}
}

type WelcomeEmailData struct {
	User User
	Data interface{}
}

type VerifyEmailData struct {
	User  User
	Token string
	URL   string
	Data  interface{}
}

type SendTFAEmailData struct {
	User User
	OTP  string
	Data interface{}
}

// EmailTemplater stores pointers to email templates.
type EmailTemplater struct {
	WelcomeTemplate       *template.Template
	ResetPasswordTemplate *template.Template
	InviteTemplate        *template.Template
	VerifyTemplate        *template.Template
	TFATemplate           *template.Template
}

// NewEmailTemplater creates new email templater.
func NewEmailTemplater(staticFilesStorage StaticFilesStorage) (*EmailTemplater, error) {
	et := EmailTemplater{}
	var err error

	if et.InviteTemplate, err = staticFilesStorage.ParseTemplate(StaticPagesNames.InviteEmail); err != nil {
		return nil, err
	}
	if et.ResetPasswordTemplate, err = staticFilesStorage.ParseTemplate(StaticPagesNames.ResetPasswordEmail); err != nil {
		return nil, err
	}
	if et.TFATemplate, err = staticFilesStorage.ParseTemplate(StaticPagesNames.TFAEmail); err != nil {
		return nil, err
	}
	if et.VerifyTemplate, err = staticFilesStorage.ParseTemplate(StaticPagesNames.VerifyEmail); err != nil {
		return nil, err
	}
	if et.WelcomeTemplate, err = staticFilesStorage.ParseTemplate(StaticPagesNames.WelcomeEmail); err != nil {
		return nil, err
	}
	return &et, nil
}
