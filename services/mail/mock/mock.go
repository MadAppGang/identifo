package mock

import (
	"fmt"
	"log/slog"

	"github.com/madappgang/identifo/v2/model"
)

type EmailService struct {
	logger       *slog.Logger
	sendMessages []string
}

// NewTransport creates new email mock transport, all it does just prints everything to console.
func NewTransport(logger *slog.Logger) model.EmailTransport {
	return &EmailService{logger: logger}
}

// SendMessage returns nil error.
func (es *EmailService) SendMessage(subject, body, recipient string) error {
	msg := fmt.Sprintf("✉️: MOCK EMAIL SERVICE: Sending message \nsubject: %s\nbody: %s\n recipient: %s\n\n", subject, body, recipient)

	es.logger.Debug("✉️: MOCK EMAIL SERVICE: Sending message",
		"subject", subject,
		"body", body,
		"recipient", recipient)

	es.sendMessages = append(es.sendMessages, msg)
	return nil
}

// SendHTML returns nil error.
func (es *EmailService) SendHTML(subject, html, recipient string) error {
	msg := fmt.Sprintf("✉️: MOCK EMAIL SERVICE: Sending HTML \nsubject: %s\nhtml: %s\n recipient: %s\n\n", subject, html, recipient)

	es.logger.Debug("✉️: MOCK EMAIL SERVICE: Sending HTML",
		"subject", subject,
		"html", html,
		"recipient", recipient)

	es.sendMessages = append(es.sendMessages, msg)
	return nil
}

func (es *EmailService) Messages() []string {
	return es.sendMessages
}
