package storage_test

import (
	"context"
	"testing"

	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage"
	"github.com/stretchr/testify/require"
)

func TestRequestSMSTest(t *testing.T) {
	cc := storage.UserStorageController{}

	ch := model.UserAuthChallenge{
		UserID:            "u1",
		DeviceID:          "d1",
		UserCodeChallenge: "ucc1234",
		OTP:               "123",
		Strategy: model.FirstFactorInternalStrategy{
			Challenge: model.AuthChallengeTypeOTP,
			Transport: model.AuthTransportTypeSMS,
		},
	}
	ch, err := cc.RequestChallenge(context.TODO(), ch)
	require.NoError(t, err)
}


