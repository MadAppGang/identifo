package mock

import "fmt"

// SMSServiceMock mocks SMS service.
type SMSServiceMock struct{}

// NewSMSService returns pointer to newly created SMS service mock.
func NewSMSService() (*SMSServiceMock, error) {
	return &SMSServiceMock{}, nil
}

// SendSMS implements SMSService.
func (ss *SMSServiceMock) SendSMS(recipient, message string) error {
	fmt.Printf("📱: MOCK SMS SERVICE: Sending SMS \recipient: %s\nmessage: %s\n\n", recipient, message)
	return nil
}
