package mailgun

import (
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
}

//NewEmailServiceFromEnv create mail service getting all settings from env
func NewEmailServiceFromEnv() model.EmailService {
	es := emailService{}
	domain := os.Getenv(MailgunDomainKey)
	privateKey := os.Getenv(MailgunPrivateKey)
	publicKey := os.Getenv(MailgunPublicKey)
	sender := os.Getenv(MAilgunSenderKey)
	mg := mailgun.NewMailgun(domain, privateKey, publicKey)
	es.mailgun = mg
	es.sender = sender
	return es
}

// NewEmailService creates and intiate email service
func NewEmailService(domain, apiKey, publicAPIKey, sender string) model.EmailService {
	es := emailService{}
	mg := mailgun.NewMailgun(domain, apiKey, publicAPIKey)
	es.mailgun = mg
	es.sender = sender
	return es
}

func (es emailService) SendMessage(subject, body, recipient string) (string, string, error) {
	message := es.mailgun.NewMessage(es.sender, subject, body, recipient)
	return es.mailgun.Send(message)
}

func (es emailService) SendHTML(subject, html, recipient string) (string, string, error) {
	message := es.mailgun.NewMessage(es.sender, subject, "", recipient)
	message.SetHtml(html)
	return es.mailgun.Send(message)
}
