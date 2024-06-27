package mock

import (
	"log/slog"
)

// SMSServiceMock mocks SMS service.
type SMSServiceMock struct {
	logger *slog.Logger
}

// NewSMSService returns pointer to newly created SMS service mock.
func NewSMSService(logger *slog.Logger) (*SMSServiceMock, error) {
	return &SMSServiceMock{logger}, nil
}

// SendSMS implements SMSService.
func (ss *SMSServiceMock) SendSMS(recipient, message string) error {
	ss.logger.Info("📱: MOCK SMS SERVICE: Sending SMS",
		"recipient", recipient,
		"message", message)
	return nil
}
