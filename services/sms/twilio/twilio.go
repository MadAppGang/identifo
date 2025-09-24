package twilio

import (
	"errors"
	"fmt"
	"log/slog"
	"regexp"

	"github.com/madappgang/identifo/v2/model"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

// SMSService sends SMS via Twilio service.
type SMSService struct {
	logger              *slog.Logger
	messagingServiceSid string
	sendFrom            string
	client              *twilio.RestClient
	noopNumbersRegexps  []*regexp.Regexp
}

// NewSMSService creates, inits and returns Twilio-backed SMS service.
func NewSMSService(
	logger *slog.Logger,
	settings model.TwilioServiceSettings,
) (*SMSService, error) {
	t := &SMSService{
		logger:              logger,
		messagingServiceSid: settings.ServiceSid,
		sendFrom:            settings.SendFrom,
		client: twilio.NewRestClientWithParams(twilio.ClientParams{
			Username: settings.AccountSid,
			Password: settings.AuthToken,
		}),
		noopNumbersRegexps: make([]*regexp.Regexp, len(settings.NoopNumbersRegexPatterns)),
	}

	if len(settings.Region) > 0 {
		t.client.Region = settings.Region
	}
	if len(settings.Edge) > 0 {
		t.client.Edge = settings.Edge
	}

	for i, r := range settings.NoopNumbersRegexPatterns {
		compiled, err := regexp.Compile(r)
		if err != nil {
			return nil, fmt.Errorf("compiling noop phone number regexp %q: %s", r, err)
		}

		t.noopNumbersRegexps[i] = compiled
	}

	return t, nil
}

// SendSMS sends SMS messages using Twilio service.
func (ss *SMSService) SendSMS(recipient, message string) error {
	if ss.client == nil {
		return errors.New("twilio SMS service is not configured")
	}

	for _, r := range ss.noopNumbersRegexps {
		if r.MatchString(recipient) {
			ss.logger.Info("Will not attempt sending SMS for a no-op number",
				"phoneNumber", recipient, "regexp", r.String())

			return nil
		}
	}

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(recipient)
	params.SetBody(message)

	// send from uses exact phone number to send from
	if len(ss.sendFrom) > 0 {
		params.SetFrom(ss.sendFrom)
	} else
	// sending messages with copilot
	// https://www.twilio.com/docs/messaging/services#send-a-message-with-copilot
	if len(ss.messagingServiceSid) > 0 {
		params.SetMessagingServiceSid(ss.messagingServiceSid)
	} else {
		return errors.New("twilio SMS service has no sendFrom nor messagingServiceSid for sending the message configured")
	}
	resp, err := ss.client.Api.CreateMessage(params)
	if err != nil {
		return err
	}

	ss.logger.Info("Twilio service sending SMS",
		"response", resp)
	return nil
}
