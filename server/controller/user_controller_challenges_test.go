package controller_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/server/controller"
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
	cc := controller.NewUserStorageController(
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
	_, err := cc.RequestChallenge(context.TODO(), ch, "+61450123456")
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

	cc := controller.NewUserStorageController(
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
	_, err := cc.RequestChallenge(context.TODO(), ch, "+61450123456")
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
				Identity:  model.AuthIdentityTypePhone,
			},
		},
	})

	// auth storage
	uas := &mock.UserAuthStorage{}

	// user storage
	u := &mock.UserStorage{}
	u.Users = append(u.Users, model.User{
		ID:          "u1",
		PhoneNumber: "+61450123456",
	})

	// sms service
	ss := &smock.SMSService{}

	cc := controller.NewUserStorageController(
		u,   // u
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
			Identity:  model.AuthIdentityTypePhone,
			Transport: model.AuthTransportTypeSMS,
		},
	}
	_, err := cc.RequestChallenge(context.TODO(), ch, "+61450123456")
	require.NoError(t, err)
	e, _ := ss.Last()
	require.Equal(t, 1, len(uas.Challenges))
	ch, _ = uas.GetLatestChallenge(context.TODO(), ch.Strategy, ch.UserID)
	assert.Equal(t, fmt.Sprintf("Here is the OTP code %s, expire in 5 minutes.", ch.OTP), e)
}

func TestRequestSMSSendUK(t *testing.T) {
	dm := map[model.SMSMessageType]string{
		model.SMSMessageTypeOTPMagicLink: "Here is your test magic link {{.URL}} which will expire in {{.Expires}} mins.",
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
				Identity:  model.AuthIdentityTypePhone,
			},
		},
	})

	// auth storage
	uas := &mock.UserAuthStorage{}

	// user storage
	u := &mock.UserStorage{}
	u.Users = append(u.Users, model.User{
		ID:          "u1",
		PhoneNumber: "+61450123456",
	})

	// sms service
	ss := &smock.SMSService{}

	cc := controller.NewUserStorageController(
		u,   // u
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
			Identity:  model.AuthIdentityTypePhone,
		},
	}

	_, err := cc.RequestChallenge(context.TODO(), ch, "+61450123456")
	require.NoError(t, err)
	e, _ := ss.Last()
	ch, _ = uas.GetLatestChallenge(context.TODO(), ch.Strategy, ch.UserID)
	// localized version should be here, as get user by strategy not implemented yet, we are getting english result here
	assert.Equal(t, fmt.Sprintf("Here is your test magic link http://localhost/web/otp_confirm?appId=a1&amp;otp=%s which will expire in 30 mins.", ch.OTP), e)

	// user with ukrainian locale
	u.Users = append(u.Users, model.User{
		ID:          "u2",
		PhoneNumber: "+61450123455",
		Locale:      language.Ukrainian.String(),
	})
	_, err = cc.RequestChallenge(context.TODO(), ch, "+61450123455")
	require.NoError(t, err)
	ch, _ = uas.GetLatestChallenge(context.TODO(), ch.Strategy, ch.UserID)
	e, _ = ss.Last()
	assert.Equal(t, fmt.Sprintf("Тримай своє посилання http://localhost/web/otp_confirm?appId=a1&amp;otp=%s яке буде недійсно через 30 хвилин.", ch.OTP), e)
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
				Identity:  model.AuthIdentityTypeEmail,
			},
		},
	})

	// auth storage
	uas := &mock.UserAuthStorage{}

	// user storage
	u := &mock.UserStorage{}
	u.Users = append(u.Users, model.User{
		ID:    "u1",
		Email: "mail@aooth.com",
	})

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

	cc := controller.NewUserStorageController(
		u,   // u
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
			Identity:  model.AuthIdentityTypeEmail,
		},
	}
	_, err := cc.RequestChallenge(context.TODO(), ch, "mail@aooth.com")
	require.NoError(t, err)

	emailTransport, ok := es.Transport().(*emock.EmailService)
	require.True(t, ok)
	// TODO: now no messages send as user has no email, that's why it shows ok, but it's not ok
	require.Equal(t, 1, len(emailTransport.Messages()))
	e, _ := emailTransport.Messages()[len(emailTransport.Messages())-1]["body"]
	ch, _ = uas.GetLatestChallenge(context.TODO(), ch.Strategy, ch.UserID)
	assert.Contains(t, e, fmt.Sprintf("Click <a href=\"http://localhost/web/otp_confirm?appId=a1&otp=%s\">here</a> to login.", ch.OTP))
}
