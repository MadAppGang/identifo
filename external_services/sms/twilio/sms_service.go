package twilio

import (
	"errors"
	"fmt"

	"github.com/madappgang/identifo/model"
	"github.com/sfreiberg/gotwilio"
)

// SMSService sends SMS via Twilio service.
type SMSService struct {
	messagingServiceSid string
	client              *gotwilio.Twilio
}

// NewSMSService creates, inits and returns Twilio-backed SMS service.
func NewSMSService(settings model.SMSServiceSettings) (*SMSService, error) {
	if len(settings.AccountSid) == 0 || len(settings.AuthToken) == 0 || len(settings.ServiceSid) == 0 {
		return nil, fmt.Errorf("Error creating Twilio SMS service, missing at least one of the parameters:"+
			"\n sidKey : %v\n tokenKey : %v\n ServiceSidKey : %v\n", settings.AccountSid, settings.AuthToken, settings.ServiceSid)
	}

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
	_, _, err := ss.client.SendSMSWithCopilot(ss.messagingServiceSid, recipient, message, "", "")
	return err
}
