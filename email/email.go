package email

import (
	"github.com/madappgang/identifo/model"
	mailgun "github.com/mailgun/mailgun-go"
)

type emailService struct {
	mailgun mailgun.Mailgun
	sender  string
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
