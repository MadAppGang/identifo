package nexmo

import (
	"errors"
	"fmt"

	"github.com/madappgang/identifo/model"
	"github.com/njern/gonexmo"
)

// SMSService sends SMS via Twilio service.
type SMSService struct {
	messagingServiceSid string
	client              *nexmo.Client
}

// NewSMSService creates, inits and returns Twilio-backed SMS service.
func NewSMSService(settings model.SMSServiceSettings) (*SMSService, error) {
	client, err := nexmo.NewClient(settings.ApiKey, settings.ApiSecret)
	if err != nil {
		return nil, err
	}
	t := &SMSService{
		messagingServiceSid: settings.ServiceSid,
		client:              client,
	}
	return t, nil
}

// SendSMS sends SMS messages using Twilio service.
func (ss *SMSService) SendSMS(recipient, message string) error {
	if ss.client == nil {
		return errors.New("Nexmo SMS service is not configured. ")
	}
	msg := &nexmo.SMSMessage{
		From: "Nexmo", // It's not used in normal mode
		To:   recipient,
		Type: nexmo.Text,
		Text: message,
	}
	resp, err := ss.client.SMS.Send(msg)
	fmt.Println(resp)
	return err
}
