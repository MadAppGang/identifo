package mailgun

import (
	"bytes"
	"html/template"

	"github.com/madappgang/identifo/model"
	mailgun "github.com/mailgun/mailgun-go"
)

type emailService struct {
	mailgun mailgun.Mailgun
	sender  string
	tmpltr  *model.EmailTemplater
}

// NewEmailService creates and inits new email service.
func NewEmailService(ess model.EmailServiceSettings, templater *model.EmailTemplater) model.EmailService {
	mg := mailgun.NewMailgun(ess.Domain, ess.PrivateKey, ess.PublicKey)
	return emailService{mailgun: mg, sender: ess.Sender, tmpltr: templater}
}

// SendMessage sends email with plain text.
func (es emailService) SendMessage(subject, body, recipient string) error {
	message := es.mailgun.NewMessage(es.sender, subject, body, recipient)
	_, _, err := es.mailgun.Send(message)
	return err
}

// SendHTML sends email with html.
func (es emailService) SendHTML(subject, html, recipient string) error {
	message := es.mailgun.NewMessage(es.sender, subject, "", recipient)
	message.SetHtml(html)
	_, _, err := es.mailgun.Send(message)
	return err
}

// Templater returns email service templater.
func (es emailService) Templater() *model.EmailTemplater {
	return es.tmpltr
}

// SendTemplateEmail applies html template to the specified data and sends it in an email.
func (es emailService) SendTemplateEmail(subject, recipient string, template *template.Template, data interface{}) error {
	var tpl bytes.Buffer
	if err := template.Execute(&tpl, data); err != nil {
		return err
	}
	return es.SendHTML(subject, tpl.String(), recipient)
}

// SendResetEmail sends reset password emails.
func (es emailService) SendResetEmail(subject, recipient string, data model.ResetEmailData) error {
	return es.SendTemplateEmail(subject, recipient, es.tmpltr.ResetPasswordTemplate, data)
}

// SendInviteEmail sends invite email to the recipient.
func (es emailService) SendInviteEmail(subject, recipient string, data model.InviteEmailData) error {
	return es.SendTemplateEmail(subject, recipient, es.tmpltr.InviteTemplate, data)
}

// SendWelcomeEmail sends welcoming emails.
func (es emailService) SendWelcomeEmail(subject, recipient string, data model.WelcomeEmailData) error {
	return es.SendTemplateEmail(subject, recipient, es.tmpltr.WelcomeTemplate, data)
}

// SendVerifyEmail sends verification emails.
func (es emailService) SendVerifyEmail(subject, recipient string, data model.VerifyEmailData) error {
	return es.SendTemplateEmail(subject, recipient, es.tmpltr.VerifyTemplate, data)
}

// SendTFAEmail sends emails with one-time password.
func (es emailService) SendTFAEmail(subject, recipient string, data model.SendTFAEmailData) error {
	return es.SendTemplateEmail(subject, recipient, es.tmpltr.TFATemplate, data)
}
