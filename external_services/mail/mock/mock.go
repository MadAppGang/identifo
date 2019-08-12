package mock

import (
	"html/template"

	"github.com/madappgang/identifo/model"
)

type emailService struct{}

// NewEmailService creates new email service mock.
func NewEmailService() model.EmailService {
	return &emailService{}
}

// SendMessage returns nil error.
func (es emailService) SendMessage(subject, body, recipient string) error {
	return nil
}

// SendHTML returns nil error.
func (es emailService) SendHTML(subject, html, recipient string) error {
	return nil
}

// Templater implements model.EmailService.
func (es emailService) Templater() *model.EmailTemplater {
	return nil
}

// SendTemplateEmail returns nil error.
func (es emailService) SendTemplateEmail(subject, recipient string, template *template.Template, data interface{}) error {
	return nil
}

// SendResetEmail returns nil error.
func (es emailService) SendResetEmail(subject, recipient string, data interface{}) error {
	return nil
}

// SendInviteEmail returns nil error.
func (es emailService) SendInviteEmail(subject, recipient string, data interface{}) error {
	return nil
}

// SendWelcomeEmail returns nil error.
func (es emailService) SendWelcomeEmail(subject, recipient string, data interface{}) error {
	return nil
}

// SendVerifyEmail returns nil error.
func (es emailService) SendVerifyEmail(subject, recipient string, data interface{}) error {
	return nil
}

// SendTFAEmail returns nil error.
func (es emailService) SendTFAEmail(subject, recipient string, data interface{}) error {
	return nil
}
