package mock

import (
	"fmt"

	"github.com/madappgang/identifo/model"
)

type emailService struct{}

// NewTransport creates new email mock transport, all it does just prints everyting to console.
func NewTransport() model.EmailTransport {
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
