package twilio_test

import (
	"testing"

	"github.com/madappgang/identifo/v2/logging"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/services/sms/twilio"
	"github.com/stretchr/testify/require"
)

func TestSensSMS_WhenNoopNumber_ReturnsNil(t *testing.T) {
	logger := logging.NewLogger("json", "info")

	settings := model.TwilioServiceSettings{
		AccountSid:               "testAccountSid",
		ServiceSid:               "testServiceSid",
		NoopNumbersRegexPatterns: []string{"\\+123456.*"},
	}

	sut, err := twilio.NewSMSService(logger, settings)
	require.NoError(t, err)

	err = sut.SendSMS("+123456789", "test message")
	require.NoError(t, err)
}
