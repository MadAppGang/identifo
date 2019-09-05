package model

import (
	"html/template"
)

// EmailService manages sending emails.
type EmailService interface {
	SendMessage(subject, body, recipient string) error
	SendHTML(subject, html, recipient string) error

	SendTemplateEmail(subject, recipient string, template *template.Template, data interface{}) error
	SendResetEmail(subject, recipient string, data interface{}) error
	SendInviteEmail(subject, recipient string, data interface{}) error
	SendWelcomeEmail(subject, recipient string, data interface{}) error
	SendVerifyEmail(subject, recipient string, data interface{}) error
	SendTFAEmail(subject, recipient string, data interface{}) error

	Templater() *EmailTemplater
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
