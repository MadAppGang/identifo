package twilio

// SMSServiceMock mocks SMS service.
type SMSServiceMock struct{}

// NewSMSServiceMock returns pointer to newly created SMS service mock.
func NewSMSServiceMock() (*SMSServiceMock, error) {
	return &SMSServiceMock{}, nil
}

// SendSMS implements SMSService.
func (ss *SMSServiceMock) SendSMS(recipient, message string) error {
	return nil
}
