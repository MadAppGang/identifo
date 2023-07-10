package mock

import "fmt"

// SMSServiceMock mocks SMS service.
type SMSService struct {
	Messages   []string
	Recipients []string
}

// NewSMSService returns pointer to newly created SMS service mock.
func NewSMSService() (*SMSService, error) {
	return &SMSService{}, nil
}

// SendSMS implements SMSService.
func (ss *SMSService) SendSMS(recipient, message string) error {
	fmt.Printf("📱: MOCK SMS SERVICE: Sending SMS: \nrecipient: %s\nmessage: %s\n\n", recipient, message)
	ss.Messages = append(ss.Messages, message)
	ss.Recipients = append(ss.Recipients, recipient)
	return nil
}

func (ss *SMSService) Last() (string, string) {
	if len(ss.Messages) == 0 {
		return "", ""
	}
	return ss.Messages[len(ss.Messages)-1], ss.Recipients[len(ss.Recipients)-1]
}
