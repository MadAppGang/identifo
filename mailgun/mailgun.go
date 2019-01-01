package mailgun

import (
	"bytes"
	"html/template"
	"os"

	"github.com/madappgang/identifo/model"
	mailgun "github.com/mailgun/mailgun-go"
)

const (
	MailgunDomainKey  = "MAILGUN_DOMAIN"
	MailgunPrivateKey = "MAILGUN_PRIVATE_KEY"
	MailgunPublicKey  = "MAILGUN_PUBLIC_KEY"
	MAilgunSenderKey  = "MAILGUN_SENDER"
)

type emailService struct {
	mailgun mailgun.Mailgun
	sender  string
	tmpltr  *model.EmailTemplater
}

//NewEmailServiceFromEnv create mail service getting all settings from env
func NewEmailServiceFromEnv(templater *model.EmailTemplater) (model.EmailService, error) {
	es := emailService{}
	domain := os.Getenv(MailgunDomainKey)
	privateKey := os.Getenv(MailgunPrivateKey)
	publicKey := os.Getenv(MailgunPublicKey)
	sender := os.Getenv(MAilgunSenderKey)
	mg := mailgun.NewMailgun(domain, privateKey, publicKey)
	es.mailgun = mg
	es.sender = sender
	if templater == nil {
		var err error
		es.tmpltr, err = model.DefaultEmailTemplater()
		if err != nil {
			return nil, err
		}
	} else {
		es.tmpltr = templater
	}
	return es, nil
}

//NewEmailService creates and inits email service
func NewEmailService(domain, apiKey, publicAPIKey, sender string, templater *model.EmailTemplater) model.EmailService {
	es := emailService{}
	mg := mailgun.NewMailgun(domain, apiKey, publicAPIKey)
	es.mailgun = mg
	es.sender = sender
	if templater != nil {
		es.tmpltr, _ = model.DefaultEmailTemplater()
	} else {
		es.tmpltr = templater
	}
	return es
}

func (es emailService) SendMessage(subject, body, recipient string) error {
	message := es.mailgun.NewMessage(es.sender, subject, body, recipient)
	_, _, err := es.mailgun.Send(message)
	return err
}

func (es emailService) SendHTML(subject, html, recipient string) error {
	message := es.mailgun.NewMessage(es.sender, subject, "", recipient)
	message.SetHtml(html)
	_, _, err := es.mailgun.Send(message)
	return err
}

//Templater returns default templater
func (es emailService) Templater() *model.EmailTemplater {
	return es.tmpltr
}

//SendTemplateEmail render data to html template and send it in email
func (es emailService) SendTemplateEmail(subject, recipient string, template *template.Template, data interface{}) error {
	var tpl bytes.Buffer
	if err := template.Execute(&tpl, data); err != nil {
		return err
	}
	return es.SendHTML(subject, tpl.String(), recipient)
}

//SendResetEmail sends reset passwords email
func (es emailService) SendResetEmail(subject, recipient string, data interface{}) error {
	return es.SendTemplateEmail(subject, recipient, es.tmpltr.ResetPasswordTemplate, data)
}

//SendWelcomeEmail sends welcome email
func (es emailService) SendWelcomeEmail(subject, recipient string, data interface{}) error {
	return es.SendTemplateEmail(subject, recipient, es.tmpltr.WelcomeTemplate, data)
}

//SendVerifyEmail sends verify email address email
func (es emailService) SendVerifyEmail(subject, recipient string, data interface{}) error {
	return es.SendTemplateEmail(subject, recipient, es.tmpltr.VerifyEmailTemplate, data)
}
