package mock

import (
	"fmt"

	"github.com/madappgang/identifo/v2/model"
)

type EmailService struct {
	SendMessages []map[string]string
}

// NewTransport creates new email mock transport, all it does just prints everything to console.
func NewTransport() model.EmailTransport {
	return &EmailService{}
}

// SendMessage returns nil error.
func (es *EmailService) SendMessage(subject, body, recipient string) error {
	msg := fmt.Sprintf("✉️: MOCK EMAIL SERVICE: Sending message \nsubject: %s\nbody: %s\n recipient: %s\n\n", subject, body, recipient)
	fmt.Printf(msg)
	message := map[string]string{
		"subject":   subject,
		"body":      body,
		"recipient": recipient,
	}
	es.SendMessages = append(es.SendMessages, message)
	return nil
}

// SendHTML returns nil error.
func (es *EmailService) SendHTML(subject, html, recipient string) error {
	msg := fmt.Sprintf("✉️: MOCK EMAIL SERVICE: Sending HTML \nsubject: %s\nhtml: %s\n recipient: %s\n\n", subject, html, recipient)
	fmt.Printf(msg)
	message := map[string]string{
		"subject":   subject,
		"body":      html,
		"recipient": recipient,
	}

	es.SendMessages = append(es.SendMessages, message)
	return nil
}

func (es *EmailService) Messages() []map[string]string {
	return es.SendMessages
}
