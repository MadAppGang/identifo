package mailgun

import (
	"bytes"
	"html/template"
	"os"

	"github.com/madappgang/identifo/model"
	mailgun "github.com/mailgun/mailgun-go"
)

const (
	// MailgunDomainKey is a name of env variable that contains Mailgun domain value.
	MailgunDomainKey = "MAILGUN_DOMAIN"
	// MailgunPrivateKey is a name of env variable that contains Mailgun private key value.
	MailgunPrivateKey = "MAILGUN_PRIVATE_KEY"
	// MailgunPublicKey is a name of env variable that contains Mailgun public key value.
	MailgunPublicKey = "MAILGUN_PUBLIC_KEY"
	// MailgunSenderKey is a name of env variable that contains Mailgun sender key value.
	MailgunSenderKey = "MAILGUN_SENDER"
)

type emailService struct {
	mailgun mailgun.Mailgun
	sender  string
	tmpltr  *model.EmailTemplater
}

// NewEmailServiceFromEnv creates new mail service getting all settings from env.
func NewEmailServiceFromEnv(templater *model.EmailTemplater) model.EmailService {
	domain := os.Getenv(MailgunDomainKey)
	privateKey := os.Getenv(MailgunPrivateKey)
	publicKey := os.Getenv(MailgunPublicKey)
	sender := os.Getenv(MailgunSenderKey)

	return NewEmailService(domain, privateKey, publicKey, sender, templater)
}

// NewEmailService creates and inits new email service.
func NewEmailService(domain, apiKey, publicAPIKey, sender string, templater *model.EmailTemplater) model.EmailService {
	mg := mailgun.NewMailgun(domain, apiKey, publicAPIKey)
	return emailService{mailgun: mg, sender: sender, tmpltr: templater}
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
func (es emailService) SendResetEmail(subject, recipient string, data interface{}) error {
	return es.SendTemplateEmail(subject, recipient, es.tmpltr.ResetPasswordTemplate, data)
}

// SendInviteEmail sends invite email to the recipient.
func (es emailService) SendInviteEmail(subject, recipient string, data interface{}) error {
	return es.SendTemplateEmail(subject, recipient, es.tmpltr.InviteEmailTemplate, data)
}

// SendWelcomeEmail sends welcoming emails.
func (es emailService) SendWelcomeEmail(subject, recipient string, data interface{}) error {
	return es.SendTemplateEmail(subject, recipient, es.tmpltr.WelcomeTemplate, data)
}

// SendVerifyEmail sends email address verification emails.
func (es emailService) SendVerifyEmail(subject, recipient string, data interface{}) error {
	return es.SendTemplateEmail(subject, recipient, es.tmpltr.VerifyEmailTemplate, data)
}
