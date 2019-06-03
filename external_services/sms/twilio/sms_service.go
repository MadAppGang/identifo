package twilio

import (
	"errors"
	"fmt"

	"github.com/sfreiberg/gotwilio"
)

// SMSService sends SMS via Twilio service.
type SMSService struct {
	messagingServiceSid string
	client              *gotwilio.Twilio
}

// NewSMSService creates, inits and returns Twilio-backed SMS service.
func NewSMSService(sidKey, tokenKey, serviceSidKey string) (*SMSService, error) {
	t := SMSService{}

	if len(sidKey) == 0 || len(tokenKey) == 0 || len(serviceSidKey) == 0 {
		return nil, fmt.Errorf("Error creating Twilio SMS service, missing param:"+
			"\n sidKey : %v\n tokenKey : %v\n ServiceSidKey : %v\n", sidKey, tokenKey, serviceSidKey)
	}
	t.messagingServiceSid = serviceSidKey
	t.client = gotwilio.NewTwilioClient(sidKey, tokenKey)
	return &t, nil
}

// SendSMS sends SMS messages using Twilio service.
func (ss *SMSService) SendSMS(recipient, message string) error {
	if ss.client == nil {
		return errors.New("Twilio SMS service is not configured")
	}
	_, _, err := ss.client.SendSMSWithCopilot(ss.messagingServiceSid, recipient, message, "", "")
	return err
}
