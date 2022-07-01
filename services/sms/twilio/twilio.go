package twilio

import (
	"errors"
	"fmt"

	"github.com/madappgang/identifo/v2/model"
	"github.com/sfreiberg/gotwilio"
)

// SMSService sends SMS via Twilio service.
type SMSService struct {
	messagingServiceSid string
	client              *gotwilio.Twilio
}

// NewSMSService creates, inits and returns Twilio-backed SMS service.
func NewSMSService(settings model.TwilioServiceSettings) (*SMSService, error) {
	t := &SMSService{
		messagingServiceSid: settings.ServiceSid,
		client:              gotwilio.NewTwilioClient(settings.AccountSid, settings.AuthToken),
	}
	return t, nil
}

// SendSMS sends SMS messages using Twilio service.
func (ss *SMSService) SendSMS(recipient, message string) error {
	if ss.client == nil {
		return errors.New("Twilio SMS service is not configured")
	}
	_, ex, err := ss.client.SendSMSWithCopilot(ss.messagingServiceSid, recipient, message, "", "")
	if err != nil {
		return err
	}

	if ex != nil {
		return fmt.Errorf("twillio error %v: %s", ex.Code, ex.Message)
	}

	return nil
}
