package sms_test

import (
	"fmt"
	"testing"

	"github.com/madappgang/identifo/v2/logging"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/services/sms"
	"github.com/madappgang/identifo/v2/services/sms/twilio"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTwilioSMSCreate(t *testing.T) {
	settings := model.SMSServiceSettings{
		Type: model.SMSServiceTwilio,
		Twilio: model.TwilioServiceSettings{
			AccountSid: "AC07374676f4e2fa65bfed5b1f5f06f36c",
			AuthToken:  "WRONG AUTH TOKEN",
			SendFrom:   "+18585440513",
		},
	}

	service, err := sms.NewService(logging.DefaultLogger, settings)
	require.NoError(t, err)

	assert.IsType(t, &twilio.SMSService{}, service)
	err = service.SendSMS("+61450396664", "I am test message")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Status: 401 - ApiError 20003: Authenticate (null) More info: https://www.twilio.com/docs/errors/20003")
	fmt.Printf("%+v\n", err)
}
