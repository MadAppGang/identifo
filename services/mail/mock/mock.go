package mock

import (
	"fmt"
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
	fmt.Printf("✉️: MOCK EMAIL SERVICE: Sending message \nsubject: %s\nbody: %s\n recipient: %s\n\n", subject, body, recipient)
	return nil
}

// SendHTML returns nil error.
func (es emailService) SendHTML(subject, html, recipient string) error {
	fmt.Printf("✉️: MOCK EMAIL SERVICE: Sending HTML \nsubject: %s\nhtml: %s\n recipient: %s\n\n", subject, html, recipient)
	return nil
}

// Templater implements model.EmailService.
func (es emailService) Templater() *model.EmailTemplater {
	return nil
}

// SendTemplateEmail returns nil error.
func (es emailService) SendTemplateEmail(subject, recipient string, template *template.Template, data interface{}) error {
	fmt.Printf("✉️: MOCK EMAIL SERVICE: Sending template Email \nsubject: %s\recipient: %s\n data: %+v\n\n", subject, recipient, data)
	return nil
}

// SendResetEmail returns nil error.
func (es emailService) SendResetEmail(subject, recipient string, data model.ResetEmailData) error {
	fmt.Printf("✉️: MOCK EMAIL SERVICE: Sending reset Email \nsubject: %s\recipient: %s\n data: %+v\n\n", subject, recipient, data)
	return nil
}

// SendInviteEmail returns nil error.
func (es emailService) SendInviteEmail(subject, recipient string, data model.InviteEmailData) error {
	fmt.Printf("✉️: MOCK EMAIL SERVICE: Sending invite Email \nsubject: %s\recipient: %s\n data: %+v\n\n", subject, recipient, data)
	return nil
}

// SendWelcomeEmail returns nil error.
func (es emailService) SendWelcomeEmail(subject, recipient string, data model.WelcomeEmailData) error {
	fmt.Printf("✉️: MOCK EMAIL SERVICE: Sending welcome Email \nsubject: %s\recipient: %s\n data: %+v\n\n", subject, recipient, data)
	return nil
}

// SendVerifyEmail returns nil error.
func (es emailService) SendVerifyEmail(subject, recipient string, data model.VerifyEmailData) error {
	fmt.Printf("✉️: MOCK EMAIL SERVICE: Sending verification Email \nsubject: %s\recipient: %s\n data: %+v\n\n", subject, recipient, data)
	return nil
}

// SendTFAEmail returns nil error.
func (es emailService) SendTFAEmail(subject, recipient string, data model.SendTFAEmailData) error {
	fmt.Printf("✉️: MOCK EMAIL SERVICE: Sending TFA Email \nsubject: %s\recipient: %s\n data: %+v\n\n", subject, recipient, data)
	return nil
}
