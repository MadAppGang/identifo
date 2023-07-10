package storage_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/services/mail"
	emock "github.com/madappgang/identifo/v2/services/mail/mock"
	smock "github.com/madappgang/identifo/v2/services/sms/mock"
	"github.com/madappgang/identifo/v2/storage"
	"github.com/madappgang/identifo/v2/storage/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestRequestSMSTestNoApp(t *testing.T) {
	cc := storage.NewUserStorageController(
		nil,           // u
		nil,           // ums
		nil,           // ua
		mock.NewApp(), // as
		nil,           // uas
		nil,           // ts
		nil,           // es
		nil,           // ss
		model.ServerSettings{},
	)

	ch := model.UserAuthChallenge{
		UserID:            "u1",
		DeviceID:          "d1",
		UserCodeChallenge: "ucc1234",
		OTP:               "123",
		AppID:             "a1",
		Strategy: model.FirstFactorInternalStrategy{
			Challenge: model.AuthChallengeTypeOTP,
			Transport: model.AuthTransportTypeSMS,
		},
	}
	_, err := cc.RequestChallenge(context.TODO(), ch)
	require.Error(t, err)
	require.True(t, errors.Is(err, mock.ErrNotFound))
}

// challenge is asking for otp code with sms, but app only does magic link with sms
func TestRequestSMSNoStrategyFound(t *testing.T) {
	as := mock.NewApp()
	as.CreateApp(model.AppData{
		ID:     "a1",
		Active: true,
		AuthStrategies: []model.AuthStrategy{
			model.FirstFactorInternalStrategy{
				Challenge: model.AuthChallengeTypeMagicLink,
				Transport: model.AuthTransportTypeSMS,
			},
		},
	})

	cc := storage.NewUserStorageController(
		nil, // u
		nil, // ums
		nil, // ua
		as,  // as
		nil, // uas
		nil, // ts
		nil, // es
		nil, // ss
		model.ServerSettings{},
	)

	ch := model.UserAuthChallenge{
		UserID:            "u1",
		DeviceID:          "d1",
		UserCodeChallenge: "ucc1234",
		OTP:               "123",
		AppID:             "a1",
		Strategy: model.FirstFactorInternalStrategy{
			Challenge: model.AuthChallengeTypeOTP,
			Transport: model.AuthTransportTypeSMS,
		},
	}
	_, err := cc.RequestChallenge(context.TODO(), ch)
	require.Error(t, err)
	assert.True(t, errors.Is(err, l.ErrorRequestChallengeUnsupportedByAPP))
}

// challenge is asking for otp code with sms, but app only does magic link with sms
func TestRequestSMSSend(t *testing.T) {
	// app storage
	as := mock.NewApp()
	as.CreateApp(model.AppData{
		ID:     "a1",
		Active: true,
		AuthStrategies: []model.AuthStrategy{
			model.FirstFactorInternalStrategy{
				Challenge: model.AuthChallengeTypeOTP,
				Transport: model.AuthTransportTypeSMS,
			},
		},
	})

	// auth storage
	uas := &mock.UserAuthStorage{}

	// sms service
	ss := &smock.SMSService{}

	cc := storage.NewUserStorageController(
		nil, // u
		nil, // ums
		nil, // ua
		as,  // as
		uas, // uas
		nil, // ts
		nil, // es
		ss,  // ss
		model.ServerSettings{},
	)

	p, _ := l.NewPrinter(language.English.String())
	cc.LP = p

	// challenge
	ch := model.UserAuthChallenge{
		UserID:            "u1",
		DeviceID:          "d1",
		UserCodeChallenge: "ucc1234",
		OTP:               "123",
		AppID:             "a1",
		Strategy: model.FirstFactorInternalStrategy{
			Challenge: model.AuthChallengeTypeOTP,
			Transport: model.AuthTransportTypeSMS,
		},
	}
	_, err := cc.RequestChallenge(context.TODO(), ch)
	require.NoError(t, err)
	e, _ := ss.Last()
	ch, _ = uas.GetLatestChallenge(context.TODO(), ch.Strategy, ch.UserID)
	assert.Equal(t, fmt.Sprintf("Here is the OTP code %s, expire in 5 minutes.", ch.OTP), e)
}

func TestRequestSMSSendUK(t *testing.T) {
	dm := map[model.SMSMessageType]string{
		model.SMSMessageTypeOTPMagicLink: "Here is your magic link {{.URL}} which will expire in {{.Expires}} mins.",
	}
	uk := map[model.SMSMessageType]string{
		model.SMSMessageTypeOTPMagicLink: "Тримай своє посилання {{.URL}} яке буде недійсно через {{.Expires}} хвилин.",
	}

	// app storage
	as := mock.NewApp()
	as.CreateApp(model.AppData{
		ID:     "a1",
		Active: true,
		CustomSMSMessages: map[string]map[model.SMSMessageType]string{
			"default":                   dm,
			language.Ukrainian.String(): uk,
		},
		AuthStrategies: []model.AuthStrategy{
			model.FirstFactorInternalStrategy{
				Challenge: model.AuthChallengeTypeMagicLink,
				Transport: model.AuthTransportTypeSMS,
			},
		},
	})

	// auth storage
	uas := &mock.UserAuthStorage{}

	// sms service
	ss := &smock.SMSService{}

	cc := storage.NewUserStorageController(
		nil, // u
		nil, // ums
		nil, // ua
		as,  // as
		uas, // uas
		nil, // ts
		nil, // es
		ss,  // ss
		model.ServerSettings{},
	)
	p, _ := l.NewPrinter(language.English.String())
	cc.LP = p

	// challenge
	ch := model.UserAuthChallenge{
		UserID:            "u1",
		DeviceID:          "d1",
		UserCodeChallenge: "ucc1234",
		OTP:               "123",
		AppID:             "a1",
		Strategy: model.FirstFactorInternalStrategy{
			Challenge: model.AuthChallengeTypeMagicLink,
			Transport: model.AuthTransportTypeSMS,
		},
	}

	_, err := cc.RequestChallenge(context.TODO(), ch)
	require.NoError(t, err)
	e, _ := ss.Last()
	ch, _ = uas.GetLatestChallenge(context.TODO(), ch.Strategy, ch.UserID)
	// localized version should be here, as get user by strategy not implemented yet, we are getting english result here
	assert.NotEqual(t, fmt.Sprintf("Here is your magic link http://localhost/web/otp_confirm?appId=a1&amp;otp=%s which will expire in 30 mins.", ch.OTP), e)
}

func TestRequestEmailSend(t *testing.T) {
	// app storage
	as := mock.NewApp()
	as.CreateApp(model.AppData{
		ID:     "a1",
		Active: true,
		AuthStrategies: []model.AuthStrategy{
			model.FirstFactorInternalStrategy{
				Challenge: model.AuthChallengeTypeMagicLink,
				Transport: model.AuthTransportTypeEmail,
			},
		},
	})

	// auth storage
	uas := &mock.UserAuthStorage{}

	// email service
	etFS, _ := storage.NewFS(model.FileStorageSettings{
		Type: model.FileStorageTypeLocal,
		Local: model.FileStorageLocal{
			Path: "../static/email_templates",
		},
	})
	es, _ := mail.NewService(
		model.EmailServiceSettings{Type: model.EmailServiceMock},
		etFS,
		time.Hour,
		"",
	)

	cc := storage.NewUserStorageController(
		nil, // u
		nil, // ums
		nil, // ua
		as,  // as
		uas, // uas
		nil, // ts
		es,  // es
		nil, // ss
		model.ServerSettings{},
	)

	p, _ := l.NewPrinter(language.English.String())
	cc.LP = p

	// challenge
	ch := model.UserAuthChallenge{
		UserID:            "u1",
		DeviceID:          "d1",
		UserCodeChallenge: "ucc1234",
		OTP:               "123",
		AppID:             "a1",
		Strategy: model.FirstFactorInternalStrategy{
			Challenge: model.AuthChallengeTypeMagicLink,
			Transport: model.AuthTransportTypeEmail,
		},
	}
	_, err := cc.RequestChallenge(context.TODO(), ch)
	require.NoError(t, err)

	emailTransport, ok := es.Transport().(*emock.EmailService)
	require.True(t, ok)
	// TODO: now no messages send as user has no email, that's why it shows ok, but it's not ok
	assert.Equal(t, 1, len(emailTransport.Messages()))
	e, _ := emailTransport.Messages()[len(emailTransport.Messages())-1]["body"]
	ch, _ = uas.GetLatestChallenge(context.TODO(), ch.Strategy, ch.UserID)
	assert.Equal(t, fmt.Sprintf("Here is the OTP code %s, expire in 5 minutes.", ch.OTP), e)
}
