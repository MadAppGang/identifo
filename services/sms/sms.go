package sms

import (
	"fmt"
	"log/slog"

	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/services/sms/mock"
	"github.com/madappgang/identifo/v2/services/sms/nexmo"
	"github.com/madappgang/identifo/v2/services/sms/routemobile"
	"github.com/madappgang/identifo/v2/services/sms/twilio"
)

func NewService(
	logger *slog.Logger,
	settings model.SMSServiceSettings,
) (model.SMSService, error) {
	switch settings.Type {
	case model.SMSServiceTwilio:
		return twilio.NewSMSService(logger, settings.Twilio)
	case model.SMSServiceNexmo:
		return nexmo.NewSMSService(settings.Nexmo)
	case model.SMSServiceRouteMobile:
		return routemobile.NewSMSService(settings.Routemobile)
	case model.SMSServiceMock:
		return mock.NewSMSService(logger)
	}
	return nil, fmt.Errorf("SMS service of type '%s' is not supported", settings.Type)
}
