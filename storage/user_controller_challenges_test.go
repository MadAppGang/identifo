package storage

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
	smock "github.com/madappgang/identifo/v2/services/sms/mock"
	"github.com/madappgang/identifo/v2/storage/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestRequestSMSTestNoApp(t *testing.T) {
	cc := UserStorageController{}
	cc.as = mock.NewApp()

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
	cc := UserStorageController{}
	cc.as = mock.NewApp()
	cc.as.CreateApp(model.AppData{
		ID:     "a1",
		Active: true,
		AuthStrategies: []model.AuthStrategy{
			model.FirstFactorInternalStrategy{
				Challenge: model.AuthChallengeTypeMagicLink,
				Transport: model.AuthTransportTypeSMS,
			},
		},
	})

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
	cc := UserStorageController{}

	// app storage
	cc.as = mock.NewApp()
	cc.as.CreateApp(model.AppData{
		ID:     "a1",
		Active: true,
		AuthStrategies: []model.AuthStrategy{
			model.FirstFactorInternalStrategy{
				Challenge: model.AuthChallengeTypeOTP,
				Transport: model.AuthTransportTypeSMS,
			},
		},
	})

	p, _ := l.NewPrinter(language.English.String())
	cc.LP = p

	// auth storage
	uas := &mock.UserAuthStorage{}
	cc.uas = uas

	// sms service
	ss := &smock.SMSService{}
	cc.ss = ss

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
