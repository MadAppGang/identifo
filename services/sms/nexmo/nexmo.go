package nexmo

import (
	"errors"

	"github.com/madappgang/identifo/v2/model"
	nexmo "github.com/njern/gonexmo"
)

// SMSService sends SMS via Nexmo service.
type SMSService struct {
	client *nexmo.Client
}

// NewSMSService creates, inits and returns Nexmo-backed SMS service.
func NewSMSService(settings model.NexmoServiceSettings) (*SMSService, error) {
	client, err := nexmo.NewClient(settings.APIKey, settings.APISecret)
	if err != nil {
		return nil, err
	}
	t := &SMSService{
		client: client,
	}
	return t, nil
}

// SendSMS sends SMS messages using Nexmo service.
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
	if err != nil {
		return err
	}
	for _, messageReport := range resp.Messages {
		if messageReport.Status != nexmo.ResponseSuccess {
			return errors.New(messageReport.ErrorText)
		}
	}

	return nil
}
