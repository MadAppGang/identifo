package sms

import (
	"fmt"

	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/services/sms/mock"
	"github.com/madappgang/identifo/v2/services/sms/nexmo"
	"github.com/madappgang/identifo/v2/services/sms/routemobile"
	"github.com/madappgang/identifo/v2/services/sms/twilio"
)

func NewService(settings model.SMSServiceSettings) (model.SMSService, error) {
	switch settings.Type {
	case model.SMSServiceTwilio:
		return twilio.NewSMSService(settings.Twilio)
	case model.SMSServiceNexmo:
		return nexmo.NewSMSService(settings.Nexmo)
	case model.SMSServiceRouteMobile:
		return routemobile.NewSMSService(settings.Routemobile)
	case model.SMSServiceMock:
		return mock.NewSMSService()
	}
	return nil, fmt.Errorf("SMS service of type '%s' is not supported", settings.Type)
}
