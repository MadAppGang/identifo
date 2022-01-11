package mailgun

import (
	"github.com/madappgang/identifo/v2/model"
	mailgun "github.com/mailgun/mailgun-go"
)

type emailService struct {
	mailgun mailgun.Mailgun
	sender  string
}

// NewTransport creates and inits new mailgun transport
func NewTransport(ess model.MailgunEmailServiceSettings) model.EmailTransport {
	mg := mailgun.NewMailgun(ess.Domain, ess.PrivateKey, ess.PublicKey)
	return emailService{mailgun: mg, sender: ess.Sender}
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
